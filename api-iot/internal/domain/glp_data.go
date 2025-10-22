package domain

import (
	"errors"
	"time"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/model"
)

type GLPData struct {
	SensorID  string   
	CreatedAt time.Time
	Temperature float64
	Humidity    float64
	SoilMoisture    float64
	LightLevel    float64
}

type GLPDataReturn struct {
	Ventilation bool `json:"ventilation"`
	Irrigation  bool `json:"irrigation"`
	Lighting    bool `json:"lighting"`
}

func NewGLPData(sensorID string, createdAt int64, temperature, humidity, soilMoisture, lightLevel float64) (*GLPData, error) {
	if sensorID == "" {
		return nil, errors.New("invalid sensor ID")
	}
	if createdAt == 0 {
		return nil, errors.New("invalid created at timestamp")
	}
	if temperature < 0 || humidity < 0 || soilMoisture < 0 || lightLevel < 0 {
		return nil, errors.New("invalid sensor data")
	}

	return &GLPData{
		SensorID:    sensorID,
		CreatedAt:   time.UnixMilli(createdAt),
		Temperature: temperature,
		Humidity:    humidity,
		SoilMoisture: soilMoisture,
		LightLevel:  lightLevel,
	}, nil
}

func (d *GLPData) Apply(parameters model.GLPParametersModel) GLPDataReturn {
	ventilation := false
	if parameters.TurnOnVentilation || d.Temperature >= parameters.MaxTemperature || d.Humidity >= parameters.MaxHumidity {
		ventilation = true
	}

	irrigation := false
	if parameters.TurnOnIrrigation || d.SoilMoisture <= parameters.MinSoilMoisture {
		irrigation = true
	}

	lighting := false
	if parameters.TurnOnLighting || d.LightLevel < parameters.MinLightLevel {
		lighting = true
	}

	return GLPDataReturn{
		Ventilation: ventilation,
		Irrigation:  irrigation,
		Lighting:    lighting,
	}
}

