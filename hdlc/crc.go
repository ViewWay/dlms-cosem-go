package hdlc

// CRC-16/CCITT HDLC style (ANSI C12.18)
// Using 0xFFFF as initial value, polynomial 0x1021

var crcTable [256]uint16

func init() {
	for i := 0; i < 256; i++ {
		crc := uint16(0)
		c := uint16(i) << 8
		for j := 0; j < 8; j++ {
			if (crc^c)&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
			c <<= 1
		}
		crcTable[i] = crc
	}
}

func reverseByte(b byte) byte {
	var r byte
	for i := 0; i < 8; i++ {
		r |= ((b >> i) & 1) << (7 - i)
	}
	return r
}

func reverseBytes(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = reverseByte(b)
	}
	return out
}

// CRCCCITT calculates HDLC CRC-16 CCITT.
func CRCCCITT(data []byte) []byte {
	reversed := reverseBytes(data)
	crc := uint16(0xFFFF)
	for _, c := range reversed {
		tmp := ((crc >> 8) & 0xFF) ^ uint16(c)
		crc = (crc << 8) ^ crcTable[tmp]
	}

	lsb := reverseByte(byte(crc&0x00FF)) ^ 0xFF
	msb := reverseByte(byte((crc&0xFF00)>>8)) ^ 0xFF
	return []byte{msb, lsb}
}

// VerifyFCS verifies the FCS of a frame.
func VerifyFCS(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	fcs := CRCCCITT(data[:len(data)-2])
	return fcs[0] == data[len(data)-2] && fcs[1] == data[len(data)-1]
}
