package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GlpData struct {
	ID					uuid.UUID	`json:"id" gorm:"type:uuid;primaryKey"`
	SensorID			string		`json:"sensor_id"`
	Temperature			float32		`json:"temperature"`
	Humidity			float32		`json:"humidity"`
	SoilMoisture		float32		`json:"soil_moisture"`
	LightLevel			float32		`json:"light_level"`
	CreatedAt	        time.Time 	`json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt			time.Time 	`json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *GlpData) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New();
	}
	return
}

func (GlpData) TableName() string {
	return "glp_data";
}

type GlpParameters struct {
	ID					uint		`json:"id" gorm:"type:uuid;primaryKey"`
	SensorID			string		`json:"sensor_id"`
	MaxTemperature		float32		`json:"max_temperature"`
	MaxHumidity			float32		`json:"max_humidity"`
	MinSoilMoisture		float32		`json:"min_soil_moisture"`
	MinLightLevel		float32		`json:"min_light_level"`
	TurnOnVentilation	bool		`json:"turn_on_ventilation"`
	TurnOnIrrigation	bool		`json:"turn_on_irrigation"`
	TurnOnLight			bool		`json:"turn_on_lighting"`
}

func (GlpParameters) TableName() string {
	return "glp_parameters";
}