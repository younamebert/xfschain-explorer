package crc16

import (
	"encoding/binary"
	"fmt"
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

func CRCS(data []byte) string {
	var crc uint16
	length := len(data)

	for i := 0; i < length; i++ {
		crc += ((crc >> 8) & 0xff00) ^ uint16(data[i])
	}

	return fmt.Sprintf("%04X", crc)

}
