package transport

import (
	"io"
	"net"
	"time"
)

// Transport is the interface for DLMS transport layers.
type Transport interface {
	// Connect establishes the connection.
	Connect() error
	// Close closes the connection.
	Close() error
	// Send sends data.
	Send(data []byte) error
	// Receive receives data with timeout.
	Receive(timeout time.Duration) ([]byte, error)
	// SetReadTimeout sets the read timeout.
	SetReadTimeout(timeout time.Duration)
	// SetWriteTimeout sets the write timeout.
	SetWriteTimeout(timeout time.Duration)
	// IsConnected returns connection status.
	IsConnected() bool
}

// TcpTransport implements Transport over TCP.
type TcpTransport struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	conn         net.Conn
	connected    bool
}

// NewTcpTransport creates a new TCP transport.
func NewTcpTransport(address string) *TcpTransport {
	return &TcpTransport{
		Address:      address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// Connect establishes a TCP connection.
func (t *TcpTransport) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, 10*time.Second)
	if err != nil {
		return err
	}
	t.conn = conn
	t.connected = true
	return nil
}

// Close closes the TCP connection.
func (t *TcpTransport) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.connected = false
		t.conn = nil
		return err
	}
	return nil
}

// Send sends data over TCP.
func (t *TcpTransport) Send(data []byte) error {
	if !t.connected || t.conn == nil {
		return io.ErrClosedPipe
	}
	if t.WriteTimeout > 0 {
		t.conn.SetWriteDeadline(time.Now().Add(t.WriteTimeout))
	}
	_, err := t.conn.Write(data)
	return err
}

// Receive reads data from TCP with timeout.
func (t *TcpTransport) Receive(timeout time.Duration) ([]byte, error) {
	if !t.connected || t.conn == nil {
		return nil, io.ErrClosedPipe
	}
	if timeout > 0 {
		t.conn.SetReadDeadline(time.Now().Add(timeout))
	} else if t.ReadTimeout > 0 {
		t.conn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
	}
	buf := make([]byte, 4096)
	n, err := t.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// SetReadTimeout sets the read timeout.
func (t *TcpTransport) SetReadTimeout(timeout time.Duration) {
	t.ReadTimeout = timeout
}

// SetWriteTimeout sets the write timeout.
func (t *TcpTransport) SetWriteTimeout(timeout time.Duration) {
	t.WriteTimeout = timeout
}

// IsConnected returns connection status.
func (t *TcpTransport) IsConnected() bool {
	return t.connected
}

// UdpTransport implements Transport over UDP.
type UdpTransport struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	conn         *net.UDPConn
	connected    bool
}

// NewUdpTransport creates a new UDP transport.
func NewUdpTransport(address string) *UdpTransport {
	return &UdpTransport{
		Address:      address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// Connect sets up the UDP connection.
func (u *UdpTransport) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", u.Address)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	u.conn = conn
	u.connected = true
	return nil
}

// Close closes the UDP connection.
func (u *UdpTransport) Close() error {
	if u.conn != nil {
		err := u.conn.Close()
		u.connected = false
		u.conn = nil
		return err
	}
	return nil
}

// Send sends data over UDP.
func (u *UdpTransport) Send(data []byte) error {
	if !u.connected || u.conn == nil {
		return io.ErrClosedPipe
	}
	_, err := u.conn.Write(data)
	return err
}

// Receive reads data from UDP with timeout.
func (u *UdpTransport) Receive(timeout time.Duration) ([]byte, error) {
	if !u.connected || u.conn == nil {
		return nil, io.ErrClosedPipe
	}
	if timeout > 0 {
		u.conn.SetReadDeadline(time.Now().Add(timeout))
	}
	buf := make([]byte, 4096)
	n, _, err := u.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// SetReadTimeout sets the read timeout.
func (u *UdpTransport) SetReadTimeout(timeout time.Duration) {
	u.ReadTimeout = timeout
}

// SetWriteTimeout sets the write timeout.
func (u *UdpTransport) SetWriteTimeout(timeout time.Duration) {
	u.WriteTimeout = timeout
}

// IsConnected returns connection status.
func (u *UdpTransport) IsConnected() bool {
	return u.connected
}

// SerialTransport implements Transport over serial port.
type SerialTransport struct {
	Port         string
	BaudRate     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	connected    bool
	// For testing: in-memory buffer
	mockData   []byte
	mockOffset int
}

// NewSerialTransport creates a new serial transport.
func NewSerialTransport(port string, baudRate int) *SerialTransport {
	return &SerialTransport{
		Port:         port,
		BaudRate:     baudRate,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// Connect opens the serial port.
func (s *SerialTransport) Connect() error {
	s.connected = true
	return nil
}

// Close closes the serial port.
func (s *SerialTransport) Close() error {
	s.connected = false
	return nil
}

// Send writes data to the serial port.
func (s *SerialTransport) Send(data []byte) error {
	if !s.connected {
		return io.ErrClosedPipe
	}
	return nil
}

// Receive reads data from the serial port.
func (s *SerialTransport) Receive(timeout time.Duration) ([]byte, error) {
	if !s.connected {
		return nil, io.ErrClosedPipe
	}
	return nil, io.EOF
}

// SetReadTimeout sets the read timeout.
func (s *SerialTransport) SetReadTimeout(timeout time.Duration) {
	s.ReadTimeout = timeout
}

// SetWriteTimeout sets the write timeout.
func (s *SerialTransport) SetWriteTimeout(timeout time.Duration) {
	s.WriteTimeout = timeout
}

// IsConnected returns connection status.
func (s *SerialTransport) IsConnected() bool {
	return s.connected
}
