package processor

import (
	"errors"
	"fmt"

	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/domain"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/model"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/protocol"
	"gorm.io/gorm"
)

type Processor struct {
	db *gorm.DB
}

func NewProcessor(db *gorm.DB) *Processor {
	return &Processor{db: db}
}

func (p *Processor) Start(header protocol.GLPHeader, payload protocol.GLPPayload) (any, error) {
	switch header.Method {
		case protocol.GLPMethodPost:
			return p.handlePost(header, payload)
		case protocol.GLPMethodGet:
			return p.handleGet(header)
		default:
			return nil, errors.New("unknown method")
	}
}

func (p *Processor) handlePost(header protocol.GLPHeader, payload protocol.GLPPayload) (any, error) {
	switch header.Route {
	case "sensor/data":
		domain, err := domain.NewGLPData(
			header.Identifier,
			payload.Timestamp,
			float64(payload.Temperature),
			float64(payload.Humidity),
			float64(payload.SoilHumidity),
			float64(payload.Light),
		)
		if err != nil {
			return nil, err
		}
		fmt.Println("Domain", domain)
		errDb := p.db.Create(&model.GLPDataModel{
			Id: domain.Id,
			SensorID:    domain.SensorID,
			CreatedAt:   domain.CreatedAt,
			Temperature: domain.Temperature,
			Humidity:    domain.Humidity,
			SoilMoisture: domain.SoilMoisture,
			LightLevel:  domain.LightLevel,
		})
		fmt.Println(errDb)
		var parameters model.GLPParametersModel
		result := p.db.Where("sensor_id = ?", domain.SensorID).First(&parameters)
		if result.Error != nil {
			return nil, result.Error
		}

		return domain.Apply(parameters), nil
	default:
		return nil, errors.New("unknown route")
	}
}

func (p *Processor) handleGet(header protocol.GLPHeader) (any, error) {
	switch header.Route {
	case "parameters":
		var parameters model.GLPParametersModel
		result := p.db.Where("sensor_id = ?", header.Identifier).First(&parameters)
		if result.Error != nil {
			return nil, result.Error
		}

		return domain.GLPDataReturn{
			Ventilation: parameters.TurnOnVentilation,
			Irrigation:  parameters.TurnOnIrrigation,
			Lighting:    parameters.TurnOnLighting,
		}, nil
	default:
		return nil, errors.New("unknown route")
	}
}