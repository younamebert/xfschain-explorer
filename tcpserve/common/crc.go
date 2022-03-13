package common

import (
	"encoding/binary"
)

func CRC(data []byte) []byte {
	var crc uint16
	length := len(data)

	for i := 0; i < length; i++ {
		crc += ((crc >> 8) & 0xff00) ^ uint16(data[i])
	}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, crc)
	return b
}

// func main() {

// 	//AA F3 0401 5E 01 60

// 	m_data := []byte{0xAA, 0xF3, 0x04, 0x01, 0x5E, 0x01, 0x60}
// 	// m_data := []byte{0xAA, 0xF1, 0x02, 0x02, 0x02}
// 	fmt.Printf("%s", CRC(m_data))
// }
