package connection

import (
	"net"

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

func (h *ConnectionHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	
	header, err := h.protocol.ParseHeader(conn)
	if err != nil {
		return
	}

	payload, err := h.protocol.ParsePayload(conn, header)
	if err != nil {
		return
	}

	
}