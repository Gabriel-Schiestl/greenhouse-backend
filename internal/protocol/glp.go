package protocol

import "net"

type GLPHeader struct {}

type GLPPayload struct {}

type GLP struct {
	
}

func (g GLP) Name() string {
	return "GLP"
}

func (g GLP) HeaderLen() int {
	return 8
}

func (g GLP) ParseHeader(conn net.Conn) (GLPHeader, error) {
	// Implementation for parsing GLP header
	return GLPHeader{}, nil
}

func (g GLP) ParsePayload(conn net.Conn) (GLPPayload, error) {
	// Implementation for parsing GLP payload
	return GLPPayload{}, nil
}