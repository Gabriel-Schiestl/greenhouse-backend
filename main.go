package main

import (
	"encoding/binary"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	header := make([]byte, 26)
	payload := []byte(`{"temperature":23.5,"humidity":45.0,"light":300.0,"soil_humidity":55.0,"timestamp":1625247600}`)

	// Fill header
	payloadLen := uint16(len(payload))
	binary.BigEndian.PutUint16(header[0:2], payloadLen)
	copy(header[2:10], []byte("device01"))
	copy(header[10:14], []byte("POST"))
	copy(header[14:26], []byte("sensor/data"))

	req := append(header, payload...)

	_, err = conn.Write(req)
	if err != nil {
		panic(err)
	}
}