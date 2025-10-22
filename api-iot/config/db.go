package config

import "os"

func GetDBConfig() string {

	return os.Getenv("DB_URL")
}