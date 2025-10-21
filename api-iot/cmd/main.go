package main

import (
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/processor"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/server"
)

func main() {
	handler := connection.NewConnectionHandler()
	processor := processor.NewProcessor(nil)
	server.Listen(handler, processor)
}