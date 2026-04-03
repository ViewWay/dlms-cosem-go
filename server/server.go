package server

import (
	"fmt"
	"sync"

	"github.com/ViewWay/dlms-cosem-go/asn1"
	"github.com/ViewWay/dlms-cosem-go/axdr"
	"github.com/ViewWay/dlms-cosem-go/core"
	"github.com/ViewWay/dlms-cosem-go/hdlc"
	"github.com/ViewWay/dlms-cosem-go/transport"
)

// CosemObjectHandler handles requests for a COSEM object.
type CosemObjectHandler interface {
	core.CosemObject
	GetAttribute(attrID int) (core.DlmsData, error)
	SetAttribute(attrID int, value core.DlmsData) error
	Action(methodID int, data []byte) ([]byte, error)
}

// ObjectRegistry stores registered COSEM objects.
type ObjectRegistry struct {
	mu      sync.RWMutex
	objects map[string]CosemObjectHandler // key: "classID:logicalName"
}

// NewObjectRegistry creates a new registry.
func NewObjectRegistry() *ObjectRegistry {
	return &ObjectRegistry{
		objects: make(map[string]CosemObjectHandler),
	}
}

// Register adds a COSEM object.
func (r *ObjectRegistry) Register(obj CosemObjectHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%d:%s", obj.ClassID(), obj.LogicalName().String())
	r.objects[key] = obj
}

// Find finds an object by class ID and logical name.
func (r *ObjectRegistry) Find(classID uint16, ln core.ObisCode) CosemObjectHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%d:%s", classID, ln.String())
	return r.objects[key]
}

// All returns all registered objects.
func (r *ObjectRegistry) All() []CosemObjectHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]CosemObjectHandler, 0, len(r.objects))
	for _, obj := range r.objects {
		result = append(result, obj)
	}
	return result
}

// DlmsServer is a DLMS/COSEM server implementation.
type DlmsServer struct {
	Transport      transport.Transport
	Registry       *ObjectRegistry
	UseHdlc        bool
	ClientAddress  int
	ServerAddress  int
	Authentication int
	SecuritySuite  int
	EncryptionKey  []byte
	AuthKey        []byte
	SystemTitle    []byte
	MaxPDUSize     int
}

// NewDlmsServer creates a new DLMS server.
func NewDlmsServer(tp transport.Transport) *DlmsServer {
	return &DlmsServer{
		Transport:      tp,
		Registry:       NewObjectRegistry(),
		UseHdlc:        true,
		ClientAddress:  1,
		ServerAddress:  16,
		Authentication: asn1.AuthNone,
		MaxPDUSize:     65535,
	}
}

// RegisterObject adds a COSEM object to the server.
func (s *DlmsServer) RegisterObject(obj CosemObjectHandler) {
	s.Registry.Register(obj)
}

// HandleAARQ processes an AARQ and returns an AARE.
func (s *DlmsServer) HandleAARQ(data []byte) ([]byte, error) {
	_, err := asn1.ParseAARQ(data)
	if err != nil {
		return nil, fmt.Errorf("server: parse AARQ: %w", err)
	}

	// Build initiate response
	initResp := s.buildInitiateResponse()
	userInfo := asn1.EncodeUserInformation(initResp)

	aare := &asn1.AARE{
		Result:                 asn1.AssocAccepted,
		ApplicationContextName: asn1.EncodeAppContextName(true, false),
		UserInformation:        userInfo,
	}

	return aare.Encode(), nil
}

// buildInitiateResponse builds the InitiateResponse.
func (s *DlmsServer) buildInitiateResponse() []byte {
	e := axdr.NewEncoder()
	e.WriteByte(0x01) // response bit string
	for i := 0; i < 8; i++ {
		e.WriteByte(0x00)
	}
	e.WriteByte(0x00) // quality of service
	e.WriteUint16(uint16(s.MaxPDUSize))
	return e.Bytes()
}

// HandleGet processes a GET request.
func (s *DlmsServer) HandleGet(data []byte) ([]byte, error) {
	d := axdr.NewDecoder(data)

	// Skip invoke ID byte
	_, _ = d.GetByte()

	// Get request type
	reqType, _ := d.GetByte()
	_ = reqType // always normal for now

	// Class ID (2 bytes)
	classID, err := d.GetUint16()
	if err != nil {
		return nil, fmt.Errorf("server: read class ID: %w", err)
	}

	// Logical name (6 bytes)
	lnBytes, err := d.GetBytes(6)
	if err != nil {
		return nil, fmt.Errorf("server: read logical name: %w", err)
	}
	ln, err := core.ObisFromBytes(lnBytes)
	if err != nil {
		return nil, fmt.Errorf("server: parse logical name: %w", err)
	}

	// Attribute count
	attrCount, _ := d.GetByte()

	// Attribute IDs
	for i := 0; i < int(attrCount); i++ {
		attrID, _ := d.GetByte()

		obj := s.Registry.Find(classID, ln)
		if obj == nil {
			// Return error: object undefined
			e := axdr.NewEncoder()
			e.WriteByte(0x04) // DataAccessResult: OBJECT_UNDEFINED
			return e.Bytes(), nil
		}

		value, err := obj.GetAttribute(int(attrID))
		if err != nil {
			e := axdr.NewEncoder()
			e.WriteByte(0x01) // HARDWARE_FAULT
			return e.Bytes(), nil
		}

		// Return success + data
		e := axdr.NewEncoder()
		e.WriteByte(0x00) // success
		e.WriteDlmsData(value)
		return e.Bytes(), nil
	}

	return nil, fmt.Errorf("server: no attributes requested")
}

// HandleSet processes a SET request.
func (s *DlmsServer) HandleSet(data []byte) ([]byte, error) {
	d := axdr.NewDecoder(data)
	_, _ = d.GetByte() // invoke ID
	_, _ = d.GetByte() // request type

	classID, _ := d.GetUint16()
	lnBytes, _ := d.GetBytes(6)
	ln, _ := core.ObisFromBytes(lnBytes)
	attrCount, _ := d.GetByte()

	for i := 0; i < int(attrCount); i++ {
		attrID, _ := d.GetByte()
		elem, err := d.GetDlmsData()
		if err != nil {
			return nil, err
		}

		obj := s.Registry.Find(classID, ln)
		if obj == nil {
			return []byte{0x04}, nil
		}

		err = obj.SetAttribute(int(attrID), elem)
		if err != nil {
			return []byte{0x01}, nil
		}
	}

	return []byte{0x00}, nil // success
}

// HandleAction processes an ACTION request.
func (s *DlmsServer) HandleAction(data []byte) ([]byte, error) {
	d := axdr.NewDecoder(data)
	_, _ = d.GetByte() // invoke ID
	_, _ = d.GetByte() // request type

	classID, _ := d.GetUint16()
	lnBytes, _ := d.GetBytes(6)
	ln, _ := core.ObisFromBytes(lnBytes)
	methodID, _ := d.GetByte()

	var actionData []byte
	if !d.Empty() {
		actionData = d.GetRaw()
	}

	obj := s.Registry.Find(classID, ln)
	if obj == nil {
		return []byte{0x04}, nil
	}

	result, err := obj.Action(int(methodID), actionData)
	if err != nil {
		return []byte{0x01}, nil
	}

	if result == nil {
		return []byte{0x00}, nil
	}
	return append([]byte{0x00}, result...), nil
}

// HandleRLRE processes a release request.
func (s *DlmsServer) HandleRLRE(data []byte) ([]byte, error) {
	rlre := &asn1.RLRE{Result: 0}
	return rlre.Encode(), nil
}

// ProcessAPDU processes a raw APDU and returns the response.
func (s *DlmsServer) ProcessAPDU(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("server: empty APDU")
	}

	switch data[0] {
	case asn1.TagAARQ:
		return s.HandleAARQ(data)
	case asn1.TagRLRE:
		return s.HandleRLRE(data)
	case 0xC0, 0xC1, 0xC2, 0xC3: // Confirmed service requests
		if len(data) < 2 {
			return nil, fmt.Errorf("server: APDU too short")
		}
		switch data[1] {
		case 0x01: // GetRequestNormal
			return s.HandleGet(data[2:])
		case 0x02: // SetRequestNormal
			return s.HandleSet(data[2:])
		case 0x03: // ActionRequestNormal
			return s.HandleAction(data[2:])
		default:
			return nil, fmt.Errorf("server: unsupported confirmed service 0x%02x", data[1])
		}
	default:
		return nil, fmt.Errorf("server: unsupported APDU tag 0x%02x", data[0])
	}
}

// HandleHdlcFrame processes an HDLC frame and returns the response frame.
func (s *DlmsServer) HandleHdlcFrame(frameData []byte) ([]byte, error) {
	frame, err := hdlc.ParseFrame(frameData)
	if err != nil {
		return nil, fmt.Errorf("server: HDLC parse: %w", err)
	}

	// Check control byte for UA/SNRM
	switch {
	case hdlc.IsSNRM(frame.Control):
		// Send UA response
		return hdlc.EncodeFrame(
			frame.SrcAddr, frame.DestAddr,
			hdlc.UAControl(true),
			nil,
		), nil
	case hdlc.IsDISC(frame.Control):
		return hdlc.EncodeFrame(
			frame.SrcAddr, frame.DestAddr,
			hdlc.UAControl(true),
			nil,
		), nil
	}

	if frame.Info == nil {
		return nil, nil
	}

	// Process APDU
	resp, err := s.ProcessAPDU(frame.Info)
	if err != nil {
		return nil, err
	}

	return hdlc.EncodeFrame(
		frame.SrcAddr, frame.DestAddr,
		hdlc.IFrameControl(0, 0),
		resp,
	), nil
}

// Serve starts handling requests from the transport.
func (s *DlmsServer) Serve() error {
	if err := s.Transport.Connect(); err != nil {
		return fmt.Errorf("server connect: %w", err)
	}
	defer s.Transport.Close()

	for {
		data, err := s.Transport.Receive(0)
		if err != nil {
			return err
		}

		var resp []byte
		if s.UseHdlc {
			parser := hdlc.NewFrameParser()
			frames := parser.Feed(data)
			for _, frameData := range frames {
				resp, err = s.HandleHdlcFrame(frameData)
				if err != nil {
					continue
				}
				if resp != nil {
					s.Transport.Send(resp)
				}
			}
		} else {
			resp, err = s.ProcessAPDU(data)
			if err != nil {
				continue
			}
			if resp != nil {
				s.Transport.Send(resp)
			}
		}
	}
}
