package connection

import (
	"encoding/json"
	"net"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/processor"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/protocol"
)


var Handler *ConnectionHandler

type ConnectionHandler struct {
	protocol protocol.Protocol[protocol.GLPHeader, protocol.GLPPayload]
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		protocol: protocol.GLP{},
	}
}

func (h *ConnectionHandler) HandleConnection(conn net.Conn, processor *processor.Processor) {
	defer conn.Close()
	
	header, err := h.protocol.ParseHeader(conn)
	if err != nil {
		errResponse := h.buildErrorResponse(err)

		conn.Write(errResponse)
		return
	}

	payload, err := h.protocol.ParsePayload(conn, header)
	if err != nil {
		errResponse := h.buildErrorResponse(err)

		conn.Write(errResponse)
		return
	}

	result := processor.Start(header, payload)

	response, err := json.Marshal(result)
	if err != nil {
		errResponse := h.buildErrorResponse(err)

		conn.Write(errResponse)
		return
	}

	conn.Write(response)
}

func (h *ConnectionHandler) buildErrorResponse(err error) []byte {
	result := map[string]any{"error": err.Error()}
	jsonResponse, _ := json.Marshal(result)

	return jsonResponse
}