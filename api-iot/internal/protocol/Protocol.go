package protocol

import "io"

type Protocol[H, P any] interface {
	Name() string
	HeaderLen() int
	ParseHeader(reader io.Reader) (H, error)
	ParsePayload(reader io.Reader, header H) (P, error)
	BuildResponse(data any) ([]byte, error)
	BuildErrorResponse(err error) []byte
}