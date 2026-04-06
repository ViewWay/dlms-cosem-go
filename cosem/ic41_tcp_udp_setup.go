package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// TCPUDPSetup is the COSEM TCP-UDP Setup interface class (IC 41).
type TCPUDPSetup struct {
	LogicalName        core.ObisCode
	TCPConnections     []TCPUDPSetupConnection
	UDPConnections     []TCPUDPSetupConnection
	MaximumSimultaneous uint8
	Version            uint8
}

// TCPUDPSetupConnection represents a TCP/UDP connection configuration.
type TCPUDPSetupConnection struct {
	RemoteAddress [4]byte
	RemotePort    uint16
	LocalPort     uint16
}

func (s *TCPUDPSetup) ClassID() uint16               { return core.ClassIDIPv4TCPSetup }
func (s *TCPUDPSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *TCPUDPSetup) GetVersion() uint8             { return s.Version }
func (s *TCPUDPSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *TCPUDPSetup) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, c := range s.TCPConnections {
		arr = append(arr, core.StructureData{
			core.OctetStringData(c.RemoteAddress[:]),
			core.UnsignedLongData(c.RemotePort),
			core.UnsignedLongData(c.LocalPort),
		})
	}
	return arr.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TCPUDPSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("tcp_udp_setup: expected array")
	}
	s.TCPConnections = make([]TCPUDPSetupConnection, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("tcp_udp_setup: invalid entry at index %d", i)
		}
		if addr, ok := st[0].(core.OctetStringData); ok && len(addr) == 4 {
			s.TCPConnections[i].RemoteAddress = [4]byte{addr[0], addr[1], addr[2], addr[3]}
		}
		s.TCPConnections[i].RemotePort = uint16(st[1].(core.UnsignedLongData))
		s.TCPConnections[i].LocalPort = uint16(st[2].(core.UnsignedLongData))
	}
	return nil
}
