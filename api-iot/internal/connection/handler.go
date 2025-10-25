package connection

import (
	"fmt"
	"net"
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/utils"
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
	if tcp, ok := conn.(*net.TCPConn); ok {
        tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(30 * time.Second)
    }
	
	for {
		header, err := h.protocol.ParseHeader(conn)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to parse header")

			errResponse := h.protocol.BuildErrorResponse(err)

			conn.Write(errResponse)
			break
		}
		
		utils.Logger.Info().Interface("header", header).Msg("Received header")

		var payload protocol.GLPPayload
		if header.Method == protocol.GLPMethodGet {
			payload = protocol.GLPPayload{}
		} else {
			payload, err = h.protocol.ParsePayload(conn, header)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to parse payload")

				errResponse := h.protocol.BuildErrorResponse(err)

				conn.Write(errResponse)
				break
			}
		}
		
		utils.Logger.Info().Interface("payload", payload).Msg("Received payload")

		result, processErr := processor.Start(header, payload)
		if processErr != nil {
			utils.Logger.Error().Err(processErr).Msg("Processing error")

			errResponse := h.protocol.BuildErrorResponse(processErr)

			conn.Write(errResponse)
			break
		}

		utils.Logger.Info().Interface("result", result).Msg("Processing result")

		response, err := h.protocol.BuildResponse(result)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to build response")
			
			errResponse := h.protocol.BuildErrorResponse(err)

			conn.Write(errResponse)
			break
		}

		utils.Logger.Info().Msg(fmt.Sprintf("Sending response: % X", response))

		_, writeErr := conn.Write(response)
		if writeErr != nil {
			utils.Logger.Error().Err(writeErr).Msg("Failed to send response")
			break
		}
	}
	conn.Close()
}

