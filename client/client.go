package client

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/ViewWay/dlms-cosem-go/asn1"
	"github.com/ViewWay/dlms-cosem-go/axdr"
	"github.com/ViewWay/dlms-cosem-go/core"
	"github.com/ViewWay/dlms-cosem-go/hdlc"
	"github.com/ViewWay/dlms-cosem-go/security"
	"github.com/ViewWay/dlms-cosem-go/transport"
)

// DlmsClient is a DLMS/COSEM client for communicating with meters.
type DlmsClient struct {
	Transport           transport.Transport
	UseHdlc             bool
	ClientAddress       int
	ServerAddress       int
	Authentication      int
	Password            []byte
	SecuritySuite       int
	EncryptionKey       []byte
	AuthenticationKey   []byte
	SystemTitle         []byte
	MaxPDUSize          int
	Timeout             time.Duration
	InvokeID            int
	ClientInvokeCounter uint32
	ServerInvokeCounter uint32
}

// NewDlmsClient creates a new DLMS client.
func NewDlmsClient(tp transport.Transport) *DlmsClient {
	return &DlmsClient{
		Transport:      tp,
		UseHdlc:        true,
		ClientAddress:  16,
		ServerAddress:  1,
		Authentication: asn1.AuthNone,
		MaxPDUSize:     65535,
		Timeout:        10 * time.Second,
	}
}

// nextInvokeID returns and increments the invoke ID.
func (c *DlmsClient) nextInvokeID() byte {
	id := byte(c.InvokeID & 0x0F)
	c.InvokeID = (c.InvokeID + 1) % 16
	return id
}

// Connect establishes the transport connection.
func (c *DlmsClient) Connect() error {
	return c.Transport.Connect()
}

// Disconnect closes the transport connection.
func (c *DlmsClient) Disconnect() error {
	if c.UseHdlc {
		err := c.sendDisconnect()
		if err != nil {
			// Best effort
		}
	}
	return c.Transport.Close()
}

// sendDisconnect sends DISC frame.
func (c *DlmsClient) sendDisconnect() error {
	control := hdlc.DISCControl()
	dest := &hdlc.HdlcAddress{Logical: c.ServerAddress}
	src := &hdlc.HdlcAddress{Logical: c.ClientAddress}
	frame := hdlc.EncodeFrame(dest, src, control, nil)
	return c.Transport.Send(frame)
}

// Associate sends AARQ and parses AARE response.
func (c *DlmsClient) Associate() error {
	// Build initiate request
	initReq := c.buildInitiateRequest()
	userInfo := asn1.EncodeUserInformation(initReq)

	aarq := &asn1.AARQ{
		ApplicationContextName: asn1.EncodeAppContextName(true, c.SecuritySuite > 0),
		CallingAPTitle:         c.SystemTitle,
		Authentication:         c.Authentication,
		AuthenticationValue:    c.Password,
		UserInformation:        userInfo,
	}

	aarqBytes := aarq.Encode()

	// Send
	if err := c.send(aarqBytes); err != nil {
		return fmt.Errorf("associate send: %w", err)
	}

	// Receive response
	resp, err := c.receive()
	if err != nil {
		return fmt.Errorf("associate receive: %w", err)
	}

	// Parse AARE
	aare, err := asn1.ParseAARE(resp)
	if err != nil {
		return fmt.Errorf("associate parse: %w", err)
	}

	if aare.Result != asn1.AssocAccepted {
		return fmt.Errorf("association rejected: result=%d, diagnostic=%v", aare.Result, aare.ResultSourceDiagnostic)
	}

	// Parse user information (initiate response)
	if aare.UserInformation != nil {
		initResp, err := asn1.DecodeUserInformation(aare.UserInformation)
		if err == nil {
			c.parseInitiateResponse(initResp)
		}
	}

	return nil
}

// buildInitiateRequest builds the xDLMS InitiateRequest.
func (c *DlmsClient) buildInitiateRequest() []byte {
	e := axdr.NewEncoder()
	e.WriteByte(0x01) // propose-conformance: bit string
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)
	e.WriteByte(0x00)                   // proposed-quality-of-service
	e.WriteUint16(uint16(c.MaxPDUSize)) // proposed-max-pdu-size

	return e.Bytes()
}

// parseInitiateResponse parses the InitiateResponse from the server.
func (c *DlmsClient) parseInitiateResponse(data []byte) {
	d := axdr.NewDecoder(data)
	// Skip response bit string
	d.GetByte() // 0x01
	for i := 0; i < 8; i++ {
		d.GetByte()
	}
	d.GetByte() // quality of service
	if v, err := d.GetUint16(); err == nil {
		c.MaxPDUSize = int(v)
	}
}

// Get reads a COSEM attribute.
func (c *DlmsClient) Get(classID uint16, logicalName core.ObisCode, attributeID int) (core.DlmsData, error) {
	e := axdr.NewEncoder()
	invokeID := c.nextInvokeID()
	e.WriteByte(0xC0 | invokeID) // confirmed request with high priority
	e.WriteByte(0x01)            // GetRequestNormal
	e.WriteByte(byte(classID >> 8))
	e.WriteByte(byte(classID))
	e.WriteBytes(logicalName.Bytes())
	e.WriteByte(0x01) // element count
	e.WriteByte(byte(attributeID))

	reqData := e.Bytes()

	// Wrap in xDLMS APDU
	apdu := make([]byte, 2+len(reqData))
	apdu[0] = 0xC0
	apdu[1] = 0x01 // normal get
	copy(apdu[2:], reqData)

	resp, err := c.sendAndReceive(reqData)
	if err != nil {
		return nil, err
	}

	// Parse response: first byte is response type, then invoke ID, then result
	if len(resp) < 2 {
		return nil, fmt.Errorf("get: response too short")
	}
	// Skip response header
	d := axdr.NewDecoder(resp[2:])
	result, err := d.GetByte()
	if err != nil {
		return nil, fmt.Errorf("get: parse result: %w", err)
	}
	if result != 0 {
		return nil, fmt.Errorf("get: access error %d", result)
	}

	// Parse data
	elem, _, err := core.DlmsDataFromBytes(d.GetRaw())
	if err != nil {
		return nil, fmt.Errorf("get: parse data: %w", err)
	}

	return elem, nil
}

// Set writes a COSEM attribute.
func (c *DlmsClient) Set(classID uint16, logicalName core.ObisCode, attributeID int, value core.DlmsData) error {
	e := axdr.NewEncoder()
	invokeID := c.nextInvokeID()
	e.WriteByte(0xC0 | invokeID)
	e.WriteByte(0x01) // SetRequestNormal
	e.WriteByte(byte(classID >> 8))
	e.WriteByte(byte(classID))
	e.WriteBytes(logicalName.Bytes())
	e.WriteByte(0x01) // element count
	e.WriteByte(byte(attributeID))
	e.WriteDlmsData(value)

	_, err := c.sendAndReceive(e.Bytes())
	return err
}

// Action performs a COSEM action.
func (c *DlmsClient) Action(classID uint16, logicalName core.ObisCode, methodID int, data []byte) ([]byte, error) {
	e := axdr.NewEncoder()
	invokeID := c.nextInvokeID()
	e.WriteByte(0xC0 | invokeID)
	e.WriteByte(0x01) // ActionRequestNormal
	e.WriteByte(byte(classID >> 8))
	e.WriteByte(byte(classID))
	e.WriteBytes(logicalName.Bytes())
	e.WriteByte(byte(methodID))
	if data != nil {
		e.WriteBytes(data)
	} else {
		e.WriteByte(0x00) // no data
	}

	return c.sendAndReceive(e.Bytes())
}

// Release releases the association.
func (c *DlmsClient) Release() error {
	rlre := &asn1.RLRE{Result: 0}
	rlreBytes := rlre.Encode()

	if err := c.send(rlreBytes); err != nil {
		return err
	}

	_, err := c.receive()
	return err
}

// send sends data, optionally wrapping in HDLC.
func (c *DlmsClient) send(data []byte) error {
	if c.UseHdlc {
		control := hdlc.IFrameControl(0, 0)
		dest := &hdlc.HdlcAddress{Logical: c.ServerAddress}
		src := &hdlc.HdlcAddress{Logical: c.ClientAddress}
		frame := hdlc.EncodeFrame(dest, src, control, data)
		return c.Transport.Send(frame)
	}
	return c.Transport.Send(data)
}

// receive receives data, optionally unwrapping HDLC.
func (c *DlmsClient) receive() ([]byte, error) {
	data, err := c.Transport.Receive(c.Timeout)
	if err != nil {
		return nil, err
	}

	if c.UseHdlc {
		parser := hdlc.NewFrameParser()
		frames := parser.Feed(data)
		if len(frames) == 0 {
			return nil, fmt.Errorf("client: no complete HDLC frame received")
		}
		frame, err := hdlc.ParseFrame(frames[0])
		if err != nil {
			return nil, fmt.Errorf("client: HDLC parse: %w", err)
		}
		return frame.Info, nil
	}

	return data, nil
}

// sendAndReceive sends a request and waits for response.
func (c *DlmsClient) sendAndReceive(data []byte) ([]byte, error) {
	if err := c.send(data); err != nil {
		return nil, err
	}
	return c.receive()
}

// WithHdlc enables/disables HDLC framing.
func (c *DlmsClient) WithHdlc(use bool) *DlmsClient {
	c.UseHdlc = use
	return c
}

// WithAuthentication sets authentication.
func (c *DlmsClient) WithAuthentication(auth int, password []byte) *DlmsClient {
	c.Authentication = auth
	c.Password = password
	return c
}

// WithSecurity sets security suite and keys.
func (c *DlmsClient) WithSecurity(suite int, encKey, authKey []byte, systemTitle []byte) *DlmsClient {
	c.SecuritySuite = suite
	c.EncryptionKey = encKey
	c.AuthenticationKey = authKey
	c.SystemTitle = systemTitle
	return c
}

// WithAddresses sets HDLC addresses.
func (c *DlmsClient) WithAddresses(client, server int) *DlmsClient {
	c.ClientAddress = client
	c.ServerAddress = server
	return c
}

// WithTimeout sets the request timeout.
func (c *DlmsClient) WithTimeout(timeout time.Duration) *DlmsClient {
	c.Timeout = timeout
	return c
}

// GetRaw performs a raw GET and returns raw bytes.
func (c *DlmsClient) GetRaw(classID uint16, logicalName core.ObisCode, attributeID int) ([]byte, error) {
	e := axdr.NewEncoder()
	invokeID := c.nextInvokeID()
	e.WriteByte(0xC0 | invokeID)
	e.WriteByte(0x01)
	e.WriteByte(byte(classID >> 8))
	e.WriteByte(byte(classID))
	e.WriteBytes(logicalName.Bytes())
	e.WriteByte(0x01)
	e.WriteByte(byte(attributeID))
	return c.sendAndReceive(e.Bytes())
}

// Session manages a DLMS session (connect, associate, use, release, disconnect).
func (c *DlmsClient) Session(fn func(*DlmsClient) error) error {
	if err := c.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer c.Disconnect()

	if err := c.Associate(); err != nil {
		return fmt.Errorf("associate: %w", err)
	}
	defer c.Release()

	return fn(c)
}

// SecurityProcessor returns a security processor for this client.
func (c *DlmsClient) SecurityProcessor() (*security.SecurityProcessor, error) {
	return security.NewSecurityProcessor(
		c.SecuritySuite,
		c.EncryptionKey,
		c.AuthenticationKey,
		c.SystemTitle,
	)
}

// Helper to encode class ID as 2 bytes.
func encodeClassID(classID uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, classID)
	return b
}
