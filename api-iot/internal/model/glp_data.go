package model

import "time"

type GLPDataModel struct {
	SensorID  string    `json:"sensor_id" gorm:"column:sensor_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Temperature float64   `json:"temperature" gorm:"column:temperature"`
	Humidity    float64   `json:"humidity" gorm:"column:humidity"`
	SoilMoisture    float64   `json:"soil_moisture" gorm:"column:soil_moisture"`
	LightLevel    float64   `json:"light_level" gorm:"column:light_level"`
}

func (GLPDataModel) TableName() string {
	return "glp_data"
}