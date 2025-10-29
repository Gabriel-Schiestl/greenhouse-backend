package database

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() (*gorm.DB, error) {

	err := godotenv.Load(); if err != nil {
		fmt.Printf("Error loading .env file: %w", err);
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSLMODE"),
	);

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err);
	}

	fmt.Println("Connected to database successfully!")

	err = db.AutoMigrate(
		&GlpData{},
		&GlpParameters{},
		); if err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %w", err);
	}

	fmt.Println("AutoMigrate completed successfully!")

	// GlpData
	if db.Migrator().HasTable(&GlpData{}) {
        fmt.Printf("Table %s exists!", GlpData{}.TableName())
    } else {
        fmt.Println("WARNING: Table 'sensor_data' was not created!")
    }

	// GlpParameters
	if db.Migrator().HasTable(&GlpParameters{}) {
        fmt.Printf("Table %s exists!", GlpParameters{}.TableName())
    } else {
        fmt.Println("WARNING: Table 'sensor_data' was not created!")
    }

	return db, nil;
}


// Glp Parameters Fun
func GetAllParameters(db *gorm.DB) ([]GlpParameters, error) {

	var alldata []GlpParameters;
	result := db.Find(&alldata);
	return alldata, result.Error;
}

func UpdateParameter(db *gorm.DB, id uint, data map[string]interface{}) error {

	result := db.Model(&GlpParameters{}).Where("id = ?", id).Updates(data);
	if result.Error !=  nil {
		return fmt.Errorf("Failed to update specific field: %w", result.Error)
	}

	fmt.Println("Data updated successfully");
	return nil;
}


// Glp Data Fun
func GetAllData(db *gorm.DB) ([]GlpData, error) {

	var alldata []GlpData;
	result := db.Find(&alldata);
	return alldata, result.Error;
}

func AddData(db *gorm.DB, data GlpData) error {

	result := db.Create(&data); if result.Error != nil {
		return fmt.Errorf("Failed to add data: %w", result.Error)
	}

	fmt.Println("Data added successfully");
	return nil;
}

func UpdateData(db *gorm.DB, id uuid.UUID, data GlpData) error {

	// Talvez seja melhor usar map[string]interface{}

	result := db.Model(&GlpData{}).Where("id = ?", id).Updates(data);
	if result.Error !=  nil {
		return fmt.Errorf("Failed to update specific field: %w", result.Error)
	}

	fmt.Println("Data updated successfully");
	return nil;
}

func DeleteData(db *gorm.DB, id uuid.UUID) error {
	
	result := db.Delete(&GlpData{}, id);
	if result == nil {
		return fmt.Errorf("Failed to delete the data: delete returned nil result");
	}
	if result.Error != nil {
		return fmt.Errorf("Failed to delete the data: %w", result.Error);
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("No records found with id %s", id.String());
	}

	fmt.Printf("Data deleted successfully, %d row(s) affected", result.RowsAffected);
	return nil;
}