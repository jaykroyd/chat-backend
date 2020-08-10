package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// WRITE DATA
func writeLength(packet *[]byte) {
	length := len(*packet)
	writer := new(bytes.Buffer)
	if err := binary.Write(writer, binary.LittleEndian, uint32(length)); err != nil {
		panic(err)
	}

	bytes := writer.Bytes()
	// fmt.Println("Length:", bytes)

	*packet = append(bytes, *packet...)
}

func writeID(id int, packet *[]byte) {
	writer := new(bytes.Buffer)
	if err := binary.Write(writer, binary.LittleEndian, uint32(id)); err != nil {
		panic(err)
	}

	bytes := writer.Bytes()
	// fmt.Println("ID:", bytes)

	*packet = append(bytes, *packet...)
}

func writeInt(buffer *[]byte, value int) {
	writer := new(bytes.Buffer)
	if err := binary.Write(writer, binary.LittleEndian, uint32(value)); err != nil {
		panic(err)
	}

	bytes := writer.Bytes()
	// fmt.Println("WriteInt:", bytes)

	*buffer = append(*buffer, bytes...)
}

func writeBool(buffer *[]byte, value bool) {
	writer := new(bytes.Buffer)
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		panic(err)
	}

	bytes := writer.Bytes()
	// fmt.Println("WriteBool:", bytes)

	*buffer = append(*buffer, bytes...)
}

func writeFloat(buffer *[]byte, value interface{}) {
	writer := new(bytes.Buffer)
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		panic(err)
	}

	bytes := writer.Bytes()
	// fmt.Println("WriteFloat:", bytes)

	*buffer = append(*buffer, bytes...)
}

func writeString(buffer *[]byte, value string) {
	writeInt(buffer, len(value))
	bytes := []byte(value)
	// fmt.Println("WriteString:", bytes)

	*buffer = append(*buffer, bytes...)
}

// READ DATA
func unwrapPacket(packet []byte, readPos *int) (int, int) {
	size := readInt(packet, readPos)
	id := readInt(packet, readPos)
	return size, id
}

func readInt(packet []byte, readPos *int) int {
	reader := bytes.NewReader(packet[*readPos : *readPos+4])
	var i int32 = -1
	err := binary.Read(reader, binary.LittleEndian, &i)
	if err != nil {
		fmt.Println("Error reading integer at pos:", *readPos, err)
	}

	*readPos += 4
	return int(i)
}

func readBool(packet []byte, readPos *int) (bool, error) {
	reader := bytes.NewReader(packet[*readPos : *readPos+4])
	var b bool = false
	err := binary.Read(reader, binary.LittleEndian, &b)
	if err != nil {
		fmt.Println("Error reading boolean at pos:", *readPos, err)
	}

	*readPos += 4
	return b, err
}
func readFloat(packet []byte, readPos *int) (float32, error) {
	reader := bytes.NewReader(packet[*readPos : *readPos+4])
	var f float32 = -1
	err := binary.Read(reader, binary.LittleEndian, &f)
	if err != nil {
		fmt.Println("Error reading float at pos:", *readPos, err)
	}

	*readPos += 4
	return f, err
}

func readString(packet []byte, readPos *int) string {
	// Read String Length
	stringSize := readInt(packet, readPos)

	// Read String Value
	stringBytes := packet[*readPos : *readPos+stringSize]
	s := string(stringBytes)

	*readPos += stringSize
	return s
}
