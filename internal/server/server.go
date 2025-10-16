package server

import (
	"net"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
)

func Listen() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go connection.HandleConnection(conn)
	}
}