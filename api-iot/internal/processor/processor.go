package processor

import (
	"errors"

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

		go p.db.Create(&model.GLPDataModel{
			SensorID:    domain.SensorID,
			CreatedAt:   domain.CreatedAt,
			Temperature: domain.Temperature,
			Humidity:    domain.Humidity,
			SoilMoisture: domain.SoilMoisture,
			LightLevel:  domain.LightLevel,
		})

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
	case "sensor/parameters":
		var parameters model.GLPParametersModel
		result := p.db.Where("sensor_id = ?", header.Identifier).First(&parameters)
		if result.Error != nil {
			return nil, result.Error
		}

		return parameters, nil
	default:
		return nil, errors.New("unknown route")
	}
}