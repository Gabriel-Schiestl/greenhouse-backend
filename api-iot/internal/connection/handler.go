package connection

import (
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
		protocol: &protocol.GLP{},
	}
}

func (h *ConnectionHandler) HandleConnection(conn net.Conn, processor *processor.Processor) {
	for {
		header, err := h.protocol.ParseHeader(conn)
		if err != nil {
			errResponse := h.protocol.BuildErrorResponse(err)

			conn.Write(errResponse)
			break
		}

		var payload protocol.GLPPayload
		if header.Method == protocol.GLPMethodGet {
			payload = protocol.GLPPayload{}
		} else {
			payload, err = h.protocol.ParsePayload(conn, header)
			if err != nil {
				errResponse := h.protocol.BuildErrorResponse(err)

				conn.Write(errResponse)
				break
			}
		}

		result, processErr := processor.Start(header, payload)
		if processErr != nil {
			errResponse := h.protocol.BuildErrorResponse(processErr)

			conn.Write(errResponse)
			break
		}

		response, err := h.protocol.BuildResponse(result)
		if err != nil {
			errResponse := h.protocol.BuildErrorResponse(err)

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

