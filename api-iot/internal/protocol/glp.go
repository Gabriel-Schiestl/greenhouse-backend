package protocol

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type GLPMethod string

type GLPHeader struct {
	PayloadLen uint16 // 2 bytes
	Identifier string //8 bytes
	Method     GLPMethod // 4 bytes
	Route	  string // 12 bytes
}

type GLPPayload struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	Light       float32 `json:"light,omitempty"`
	SoilHumidity float32 `json:"soil_humidity,omitempty"`
	Timestamp   int64   `json:"timestamp,omitempty"`
}

type GLP struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	Light       float32 `json:"light,omitempty"`
	SoilHumidity float32 `json:"soil_humidity,omitempty"`
	Timestamp   int64   `json:"timestamp,omitempty"`
}

const (
	GLPMethodPost GLPMethod = "POST"
	GLPMethodGet  GLPMethod = "GET"
)

func (g GLP) Name() string {
	return "GLP"
}

func (g GLP) HeaderLen() int {
	return 26
}

func (g GLP) ParseHeader(conn net.Conn) (GLPHeader, error) {
	var header GLPHeader
	headerBytes := make([]byte, g.HeaderLen())
	_, err := conn.Read(headerBytes)
	if err != nil {
		return GLPHeader{}, err
	}

	header.PayloadLen = binary.BigEndian.Uint16(headerBytes[0:2])
	header.Identifier = string(headerBytes[2:10])
	header.Method = GLPMethod(headerBytes[10:14])
	header.Route = string(headerBytes[14:26])

	return header, nil
}

func (g GLP) ParsePayload(conn net.Conn, header GLPHeader) (GLPPayload, error) {
	payloadBytes := make([]byte, header.PayloadLen)
	_, err := conn.Read(payloadBytes)
	if err != nil {
		return GLPPayload{}, err
	}

	var payload GLPPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return GLPPayload{}, err
	}

	return payload, nil
}