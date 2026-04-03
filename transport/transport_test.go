package transport

import (
	"testing"
	"time"
)

func TestNewTcpTransport(t *testing.T) {
	tp := NewTcpTransport("127.0.0.1:4059")
	if tp.Address != "127.0.0.1:4059" {
		t.Error("wrong address")
	}
	if tp.ReadTimeout != 10*time.Second {
		t.Error("wrong default timeout")
	}
	if tp.IsConnected() {
		t.Error("should not be connected")
	}
}

func TestTcpTransport_Timeouts(t *testing.T) {
	tp := NewTcpTransport("127.0.0.1:4059")
	tp.SetReadTimeout(5 * time.Second)
	if tp.ReadTimeout != 5*time.Second {
		t.Error("read timeout not set")
	}
	tp.SetWriteTimeout(3 * time.Second)
	if tp.WriteTimeout != 3*time.Second {
		t.Error("write timeout not set")
	}
}

func TestTcpTransport_Send_NotConnected(t *testing.T) {
	tp := NewTcpTransport("127.0.0.1:4059")
	err := tp.Send([]byte{0x01})
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestTcpTransport_Receive_NotConnected(t *testing.T) {
	tp := NewTcpTransport("127.0.0.1:4059")
	_, err := tp.Receive(1 * time.Second)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestTcpTransport_Close_NotConnected(t *testing.T) {
	tp := NewTcpTransport("127.0.0.1:4059")
	err := tp.Close()
	if err != nil {
		t.Error("close should not error when not connected")
	}
}

func TestNewUdpTransport(t *testing.T) {
	up := NewUdpTransport("127.0.0.1:4059")
	if up.Address != "127.0.0.1:4059" {
		t.Error("wrong address")
	}
	if up.IsConnected() {
		t.Error("should not be connected")
	}
}

func TestUdpTransport_Timeouts(t *testing.T) {
	up := NewUdpTransport("127.0.0.1:4059")
	up.SetReadTimeout(5 * time.Second)
	up.SetWriteTimeout(3 * time.Second)
	if up.ReadTimeout != 5*time.Second {
		t.Error("read timeout")
	}
}

func TestUdpTransport_Send_NotConnected(t *testing.T) {
	up := NewUdpTransport("127.0.0.1:4059")
	err := up.Send([]byte{0x01})
	if err == nil {
		t.Error("expected error")
	}
}

func TestUdpTransport_Receive_NotConnected(t *testing.T) {
	up := NewUdpTransport("127.0.0.1:4059")
	_, err := up.Receive(1 * time.Second)
	if err == nil {
		t.Error("expected error")
	}
}

func TestUdpTransport_Close_NotConnected(t *testing.T) {
	up := NewUdpTransport("127.0.0.1:4059")
	err := up.Close()
	if err != nil {
		t.Error("close should not error")
	}
}

func TestNewSerialTransport(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	if sp.Port != "/dev/ttyUSB0" {
		t.Error("wrong port")
	}
	if sp.BaudRate != 9600 {
		t.Error("wrong baud rate")
	}
	if sp.IsConnected() {
		t.Error("should not be connected")
	}
}

func TestSerialTransport_ConnectClose(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	err := sp.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if !sp.IsConnected() {
		t.Error("should be connected")
	}
	err = sp.Close()
	if err != nil {
		t.Fatal(err)
	}
	if sp.IsConnected() {
		t.Error("should be disconnected")
	}
}

func TestSerialTransport_Send_NotConnected(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	err := sp.Send([]byte{0x01})
	if err == nil {
		t.Error("expected error")
	}
}

func TestSerialTransport_Receive_NotConnected(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	_, err := sp.Receive(1 * time.Second)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSerialTransport_Send_Connected(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	sp.Connect()
	err := sp.Send([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Error("send should succeed when connected")
	}
}

func TestSerialTransport_Timeouts(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	sp.SetReadTimeout(5 * time.Second)
	sp.SetWriteTimeout(3 * time.Second)
	if sp.ReadTimeout != 5*time.Second {
		t.Error("read timeout")
	}
}

func TestSerialTransport_Close_NotConnected(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	err := sp.Close()
	if err != nil {
		t.Error("close should not error")
	}
}

func TestSerialTransport_Receive_Connected(t *testing.T) {
	sp := NewSerialTransport("/dev/ttyUSB0", 9600)
	sp.Connect()
	_, err := sp.Receive(1 * time.Second)
	// Should return EOF since no real serial port
	if err == nil {
		t.Error("expected error (EOF) on mock serial")
	}
}

func TestTransport_Interface(t *testing.T) {
	var _ Transport = NewTcpTransport("127.0.0.1:4059")
	var _ Transport = NewUdpTransport("127.0.0.1:4059")
	var _ Transport = NewSerialTransport("/dev/ttyUSB0", 9600)
}
