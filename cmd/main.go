package main

import (
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/server"
)

func main() {
	handler := connection.NewConnectionHandler()
	server.Listen(handler)
}