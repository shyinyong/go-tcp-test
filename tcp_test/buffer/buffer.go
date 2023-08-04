package buffer

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func GenLoginData() *bytes.Buffer {
	//account := "m-128516svas"
	//password := "3af3cdf5a960"
	account := "m-128693affk"
	password := "d7693f7c0325"

	pkt := new(bytes.Buffer)
	binary.Write(pkt, binary.LittleEndian, uint16(1000))
	binary.Write(pkt, binary.LittleEndian, uint16(1700))
	binary.Write(pkt, binary.LittleEndian, uint8(len(account)))
	binary.Write(pkt, binary.LittleEndian, []byte(account))
	binary.Write(pkt, binary.LittleEndian, uint8(len(password)))
	binary.Write(pkt, binary.LittleEndian, []byte(password))
	binary.Write(pkt, binary.LittleEndian, uint16(19))
	binary.Write(pkt, binary.LittleEndian, int32(1))

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
	binary.Write(data, binary.LittleEndian, pkt.Bytes())

	return data
}

func GenCMDIDData(cmdID int) *bytes.Buffer {
	pkt := new(bytes.Buffer)
	binary.Write(pkt, binary.LittleEndian, cmdID)

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
	binary.Write(data, binary.LittleEndian, pkt.Bytes())

	return data
}

func StrConvToUint16(str string) uint16 {
	intNum, _ := strconv.Atoi(str)
	return uint16(intNum)
}

func BytesToUInt16(buf []byte) uint16 {
	return (uint16(buf[0]) << 8) | (uint16(buf[1]))
	// return binary.BigEndian.Uint16(buf)
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}
