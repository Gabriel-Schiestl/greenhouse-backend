package protocol

import "net"

type Protocol[H, P any] interface {
	Name() string
	HeaderLen() int
	ParseHeader(conn net.Conn) (H, error)
	ParsePayload(conn net.Conn, header H) (P, error)
}