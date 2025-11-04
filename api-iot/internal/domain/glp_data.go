package domain

import (
	"errors"
	"time"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/model"
	"github.com/google/uuid"
)

type GLPData struct {
	Id        string
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

func NewGLPData(sensorID string, temperature, humidity, soilMoisture, lightLevel float64) (*GLPData, error) {
	if sensorID == "" {
		return nil, errors.New("invalid sensor ID")
	}
	if temperature < 0 || humidity < 0 || soilMoisture < 0 || lightLevel < 0 {
		return nil, errors.New("invalid sensor data")
	}

	return &GLPData{
		Id:        uuid.New().String(),
		SensorID:    sensorID,
		CreatedAt:   time.Now(),
		Temperature: temperature,
		Humidity:    humidity,
		SoilMoisture: soilMoisture,
		LightLevel:  lightLevel,
	}, nil
}

func (d *GLPData) Apply(parameters model.GLPParametersModel) (GLPDataReturn, bool) {
	ventilation := false
	if d.Temperature >= parameters.MaxTemperature || d.Humidity >= parameters.MaxHumidity {
		ventilation = true
	}

	irrigation := false
	if d.SoilMoisture <= parameters.MinSoilMoisture {
		irrigation = true
	}

	lighting := false
	if d.LightLevel < parameters.MinLightLevel {
		lighting = true
	}

	if time.Since(parameters.LastUserUpdate) <= 1*time.Hour {
		ventilation = parameters.TurnOnVentilation
		irrigation = parameters.TurnOnIrrigation
		lighting = parameters.TurnOnLighting
	}

	return GLPDataReturn{
		Ventilation: ventilation,
		Irrigation:  irrigation,
		Lighting:    lighting,
	}, time.Since(parameters.LastUserUpdate) > 1*time.Hour
}

