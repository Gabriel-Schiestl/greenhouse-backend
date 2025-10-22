package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}