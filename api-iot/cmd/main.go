package main

import (
	"github.com/Gabriel-Schiestl/go-clarch/utils"
	"github.com/Gabriel-Schiestl/greenhouse-backend/config"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/connection"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/model"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/processor"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	dns := config.GetDBConfig()
	db, err := config.CreateDB(dns)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	db.AutoMigrate(model.GLPDataModel{}, model.GLPParametersModel{})
  
	handler := connection.NewConnectionHandler()
	processor := processor.NewProcessor(db)
	server.Listen(handler, processor)
}