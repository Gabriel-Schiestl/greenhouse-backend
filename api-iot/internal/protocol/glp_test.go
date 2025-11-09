package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGLPHeader(t *testing.T) {
	payloadLen := uint16(34)
	identifier := "DEV12345"
	method := "POST"
	route := "sensor/data "

	headerBytes := make([]byte, 26)
	binary.BigEndian.PutUint16(headerBytes[0:2], payloadLen)
	copy(headerBytes[2:10], []byte(identifier))
	copy(headerBytes[10:14], []byte(method))
	copy(headerBytes[14:26], []byte(route))

	reader := bytes.NewReader(headerBytes)

	glp := &GLP{}
	header, err := glp.ParseHeader(reader)

	require.NoError(t, err)
	require.Equal(t, payloadLen, header.PayloadLen)
	require.Equal(t, identifier, header.Identifier)
	require.Equal(t, GLPMethod(method), header.Method)
	require.Equal(t, strings.TrimSpace(route), header.Route)
}

func TestParseGLPPayload(t *testing.T) {
	payload := map[string]float32{
		"temperature":   23.5,
		"humidity":      60.0,
		"light":         100.0,
		"soil_humidity": 45.0,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	payloadLen := uint16(len(payloadBytes))
	identifier := "DEV12345"
	method := "POST"
	route := "sensor/data "

	request := make([]byte, 26 + len(payloadBytes))
	binary.BigEndian.PutUint16(request[0:2], payloadLen)
	copy(request[2:10], []byte(identifier))
	copy(request[10:14], []byte(method))
	copy(request[14:26], []byte(route))
	copy(request[26:], payloadBytes)

	reader := bytes.NewReader(request)

	glp := &GLP{}
	header, _ := glp.ParseHeader(reader)

	payloadParsed, errPayload := glp.ParsePayload(reader, header)

	require.NoError(t, errPayload)
	require.Equal(t, payload["temperature"], payloadParsed.Temperature)
	require.Equal(t, payload["humidity"], payloadParsed.Humidity)
	require.Equal(t, payload["light"], payloadParsed.Light)
	require.Equal(t, payload["soil_humidity"], payloadParsed.SoilHumidity)
}

func TestBuildResponseHeader(t *testing.T) {
	glp := &GLP{}
	res := glp.buildResponseHeader(false)

	require.Equal(t, 3, len(res))
	require.Equal(t, byte(0x00), res[0])

	errRes := glp.buildResponseHeader(true)

	require.Equal(t, 3, len(errRes))
	require.Equal(t, byte(0x01), errRes[0])
}

func TestBuildResponse(t *testing.T) {
	glp := &GLP{}

	data := map[string]string{
		"status": "ok",
	}

	marshalledData, err := json.Marshal(data)
	require.NoError(t, err)

	expectedResponse := make([]byte, 3 + len(marshalledData))
	expectedResponse[0] = 0x00
	binary.BigEndian.PutUint16(expectedResponse[1:3], uint16(len(marshalledData)))
	copy(expectedResponse[3:], marshalledData)

	response, err := glp.BuildResponse(data)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)
}

func TestBuildErrorResponse(t *testing.T) {
	glp := &GLP{}

	err := errors.New("test error")

	data := map[string]string{
		"error": err.Error(),
	}

	marshalledData, marshalErr := json.Marshal(data)
	require.NoError(t, marshalErr)

	expectedResponse := make([]byte, 3 + len(marshalledData))
	expectedResponse[0] = 0x01
	binary.BigEndian.PutUint16(expectedResponse[1:3], uint16(len(marshalledData)))
	copy(expectedResponse[3:], marshalledData)

	response := glp.BuildErrorResponse(err)
	require.Equal(t, expectedResponse, response)
}