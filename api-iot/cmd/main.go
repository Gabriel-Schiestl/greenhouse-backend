package main

import (
	"log"

	"github.com/Gabriel-Schiestl/greenhouse-backend/config"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/processor"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dns := config.GetDBConfig()
	db, err := config.CreateDB(dns)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
  
	handler := connection.NewConnectionHandler()
	processor := processor.NewProcessor(db)
	server.Listen(handler, processor)
}