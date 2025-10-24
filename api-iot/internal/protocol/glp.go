package protocol

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"
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

func (g *GLP) Name() string {
	return "GLP"
}

func (g *GLP) HeaderLen() int {
	return 26
}

func (g *GLP) ParseHeader(conn net.Conn) (GLPHeader, error) {
	if tcp, ok := conn.(*net.TCPConn); ok {
        tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(30 * time.Second)
    }

	var header GLPHeader
	headerBytes := make([]byte, g.HeaderLen())
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return GLPHeader{}, err
	}

	header.PayloadLen = binary.BigEndian.Uint16(headerBytes[0:2])
	header.Identifier = strings.TrimSpace(string(headerBytes[2:10]))
	header.Method = GLPMethod(strings.TrimSpace(string(headerBytes[10:14])))
	header.Route = strings.TrimSpace(string(headerBytes[14:26]))

	return header, nil
}

func (g *GLP) ParsePayload(conn net.Conn, header GLPHeader) (GLPPayload, error) {
	payloadBytes := make([]byte, header.PayloadLen)
	_, err := io.ReadFull(conn, payloadBytes)
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

func (g *GLP) BuildResponse(data any) ([]byte, error) {
	header := g.buildResponseHeader(false)

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(header[1:3], uint16(len(jsonResponse)))

	response := append(header, jsonResponse...)

	return response, nil
}

func (g *GLP) BuildErrorResponse(err error) []byte {
	header := g.buildResponseHeader(true)

	result := map[string]any{"error": err.Error()}
	jsonResponse, _ := json.Marshal(result)

	binary.BigEndian.PutUint16(header[1:3], uint16(len(jsonResponse)))

	response := append(header, jsonResponse...)

	return response
}

func (g *GLP) buildResponseHeader(err bool) []byte {
	headerBytes := make([]byte, 3)

	var method uint8
	if err {
		method = 1
	} else {
		method = 0
	}

	headerBytes[0] = method
	return headerBytes
}