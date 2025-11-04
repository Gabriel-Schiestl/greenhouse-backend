package model

import "time"

type GLPParametersModel struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`
	SensorID          string    `gorm:"index;not null"`
	MaxTemperature    float64   `gorm:"not null"`
	MaxHumidity       float64   `gorm:"not null"`
	MinSoilMoisture   float64   `gorm:"not null"`
	MinLightLevel     float64   `gorm:"not null"`
	TurnOnVentilation bool      `gorm:"not null"`
	TurnOnIrrigation  bool      `gorm:"not null"`
	TurnOnLighting    bool      `gorm:"not null"`
	LastUserUpdate    time.Time
}

func (GLPParametersModel) TableName() string {
	return "glp_parameters"
}