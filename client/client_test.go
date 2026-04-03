package client

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ViewWay/dlms-cosem-go/axdr"
	"github.com/ViewWay/dlms-cosem-go/core"
	"github.com/ViewWay/dlms-cosem-go/hdlc"
	"github.com/ViewWay/dlms-cosem-go/transport"
)

// mockTransport implements transport.Transport for testing
type mockTransport struct {
	data      [][]byte
	sentData  [][]byte
	connected bool
}

func (m *mockTransport) Connect() error { m.connected = true; return nil }
func (m *mockTransport) Close() error   { m.connected = false; return nil }
func (m *mockTransport) Send(data []byte) error {
	m.sentData = append(m.sentData, data)
	return nil
}
func (m *mockTransport) Receive(timeout time.Duration) ([]byte, error) {
	if len(m.data) > 0 {
		d := m.data[0]
		m.data = m.data[1:]
		return d, nil
	}
	return nil, nil
}
func (m *mockTransport) SetReadTimeout(timeout time.Duration)  {}
func (m *mockTransport) SetWriteTimeout(timeout time.Duration) {}
func (m *mockTransport) IsConnected() bool                     { return m.connected }

func TestNewDlmsClient(t *testing.T) {
	tp := &mockTransport{}
	c := NewDlmsClient(tp)
	if c == nil {
		t.Error("nil client")
	}
	if !c.UseHdlc {
		t.Error("HDLC should be enabled by default")
	}
	if c.ClientAddress != 16 {
		t.Errorf("client address=%d", c.ClientAddress)
	}
	if c.ServerAddress != 1 {
		t.Errorf("server address=%d", c.ServerAddress)
	}
}

func TestNextInvokeID(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	id1 := c.nextInvokeID()
	id2 := c.nextInvokeID()
	if id1 == id2 {
		t.Error("invoke IDs should increment")
	}
	if id1 != 0 {
		t.Errorf("first id=%d", id1)
	}
}

func TestWithHdlc(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	c.WithHdlc(false)
	if c.UseHdlc {
		t.Error("HDLC should be disabled")
	}
}

func TestWithAuthentication(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	c.WithAuthentication(2, []byte{0xAA})
	if c.Authentication != 2 {
		t.Error("wrong auth")
	}
}

func TestWithSecurity(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	key := make([]byte, 16)
	title := make([]byte, 8)
	c.WithSecurity(1, key, key, title)
	if c.SecuritySuite != 1 {
		t.Error("wrong suite")
	}
}

func TestWithAddresses(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	c.WithAddresses(32, 17)
	if c.ClientAddress != 32 || c.ServerAddress != 17 {
		t.Error("wrong addresses")
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	c.WithTimeout(5 * time.Second)
	if c.Timeout != 5*time.Second {
		t.Error("wrong timeout")
	}
}

func TestSend_NotConnected(t *testing.T) {
	tp := &mockTransport{}
	c := NewDlmsClient(tp)
	// HDLC is on, send will try to encode frame
	err := c.send([]byte{0x01, 0x02})
	// Transport is not connected but mock still accepts
	if err != nil {
		t.Error("mock should accept")
	}
}

func TestSend_NoHdlc(t *testing.T) {
	tp := &mockTransport{}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	err := c.send([]byte{0x01, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	if len(tp.sentData) != 1 {
		t.Error("should have sent 1 message")
	}
	if !bytes.Equal(tp.sentData[0], []byte{0x01, 0x02}) {
		t.Errorf("got %v", tp.sentData[0])
	}
}

func TestReceive_NoHdlc(t *testing.T) {
	tp := &mockTransport{data: [][]byte{{0x01, 0x02, 0x03}}}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	data, err := c.receive()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte{0x01, 0x02, 0x03}) {
		t.Errorf("got %v", data)
	}
}

func TestReceive_Hdlc(t *testing.T) {
	// Create a valid HDLC frame
	dest := &hdlc.HdlcAddress{Logical: 16}
	src := &hdlc.HdlcAddress{Logical: 1}
	info := []byte{0x01, 0x02, 0x03}
	frame := hdlc.EncodeFrame(dest, src, 0x00, info)

	tp := &mockTransport{data: [][]byte{frame}}
	c := NewDlmsClient(tp)
	c.WithAddresses(1, 16) // client=1, server=16 (reversed for parsing)
	data, err := c.receive()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, info) {
		t.Errorf("got %v", data)
	}
}

func TestBuildInitiateRequest(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	req := c.buildInitiateRequest()
	if len(req) < 10 {
		t.Errorf("request too short: %d", len(req))
	}
	if req[0] != 0x01 {
		t.Errorf("first byte=%02x", req[0])
	}
}

func TestParseInitiateResponse(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	// Build a minimal initiate response
	e := axdr.NewEncoder()
	e.WriteByte(0x01) // bit string
	for i := 0; i < 8; i++ {
		e.WriteByte(0x00)
	}
	e.WriteByte(0x00)   // quality
	e.WriteUint16(2048) // max PDU
	c.parseInitiateResponse(e.Bytes())
	if c.MaxPDUSize != 2048 {
		t.Errorf("max pdu=%d", c.MaxPDUSize)
	}
}

func TestConnect(t *testing.T) {
	tp := &mockTransport{}
	c := NewDlmsClient(tp)
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if !tp.connected {
		t.Error("transport should be connected")
	}
}

func TestDisconnect(t *testing.T) {
	tp := &mockTransport{connected: true}
	c := NewDlmsClient(tp)
	err := c.Disconnect()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRelease(t *testing.T) {
	tp := &mockTransport{
		data: [][]byte{
			// AARE response for release
		},
	}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	err := c.Release()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendAndReceive(t *testing.T) {
	tp := &mockTransport{data: [][]byte{{0x01, 0x02}}}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	resp, err := c.sendAndReceive([]byte{0xAA, 0xBB})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp, []byte{0x01, 0x02}) {
		t.Errorf("got %v", resp)
	}
}

func TestEncodeClassID(t *testing.T) {
	b := encodeClassID(1)
	if !bytes.Equal(b, []byte{0x00, 0x01}) {
		t.Errorf("got %v", b)
	}
	b = encodeClassID(256)
	if !bytes.Equal(b, []byte{0x01, 0x00}) {
		t.Errorf("got %v", b)
	}
}

func TestSession_ConnectError(t *testing.T) {
	tp := &mockFailTransport{}
	c := NewDlmsClient(tp)
	err := c.Session(func(c *DlmsClient) error { return nil })
	if err == nil {
		t.Error("expected error for connect failure")
	}
}

type mockFailTransport struct{}

func (m *mockFailTransport) Connect() error                                { return fmt.Errorf("connect failed") }
func (m *mockFailTransport) Close() error                                  { return nil }
func (m *mockFailTransport) Send(data []byte) error                        { return nil }
func (m *mockFailTransport) Receive(timeout time.Duration) ([]byte, error) { return nil, nil }
func (m *mockFailTransport) SetReadTimeout(timeout time.Duration)          {}
func (m *mockFailTransport) SetWriteTimeout(timeout time.Duration)         {}
func (m *mockFailTransport) IsConnected() bool                             { return false }

func TestClient_Interface(t *testing.T) {
	var _ transport.Transport = &mockTransport{}
}

func TestGetRaw(t *testing.T) {
	tp := &mockTransport{data: [][]byte{{0x00, 0x00, 0x01}}}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	_, err := c.GetRaw(1, core.ObisClock, 2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewDlmsClient(&mockTransport{})
	if c.MaxPDUSize != 65535 {
		t.Error("wrong default PDU size")
	}
	if c.Timeout != 10*time.Second {
		t.Error("wrong default timeout")
	}
	if c.Authentication != 0 {
		t.Error("wrong default auth")
	}
}

func TestSet(t *testing.T) {
	tp := &mockTransport{data: [][]byte{{0x01}}}
	c := NewDlmsClient(tp)
	c.WithHdlc(false)
	err := c.Set(1, core.ObisClock, 2, core.IntegerData(42))
	if err != nil {
		t.Fatal(err)
	}
}
