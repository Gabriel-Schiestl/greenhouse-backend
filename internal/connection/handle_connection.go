package connection

import "net"

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	// Handle the connection (read/write data)
}