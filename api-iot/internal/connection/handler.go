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
	for {
		header, err := h.protocol.ParseHeader(conn)
		if err != nil {
			errResponse := h.buildErrorResponse(err)

			conn.Write(errResponse)
			break
		}

		payload, err := h.protocol.ParsePayload(conn, header)
		if err != nil {
			errResponse := h.buildErrorResponse(err)

			conn.Write(errResponse)
			break
		}

		result, processErr := processor.Start(header, payload)
		if processErr != nil {
			errResponse := h.buildErrorResponse(processErr)

			conn.Write(errResponse)
			break
		}

		response, err := json.Marshal(result)
		if err != nil {
			errResponse := h.buildErrorResponse(err)

			conn.Write(errResponse)
			break
		}

		_, writeErr := conn.Write(response)
		if writeErr != nil {
			break
		}
	}
	conn.Close()
}

func (h *ConnectionHandler) buildErrorResponse(err error) []byte {
	result := map[string]any{"error": err.Error()}
	jsonResponse, _ := json.Marshal(result)

	return jsonResponse
}