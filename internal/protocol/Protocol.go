package protocol

type Protocol interface {
	Name() string
	HeaderLen() int
	ParseHeader([]byte) (int, error)
	ParsePayload([]byte) (any, error)
}