package server

import (
	"fmt"
	"net"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/processor"
)

func Listen(handler *connection.ConnectionHandler, processor *processor.Processor) {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println(conn.RemoteAddr())
		go handler.HandleConnection(conn, processor)
	}
}