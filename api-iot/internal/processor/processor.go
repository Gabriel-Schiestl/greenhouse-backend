package processor

import (
	"errors"

	"github.com/Gabriel-Schiestl/go-clarch/utils"
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
			utils.Logger.Warn().Msg("Unknown method")
			return nil, errors.New("unknown method")
	}
}

func (p *Processor) handlePost(header protocol.GLPHeader, payload protocol.GLPPayload) (any, error) {
	utils.Logger.Info().Msg("Handling POST request")

	switch header.Route {
	case "sensor/data":
		utils.Logger.Info().Msg("Processing sensor/data")

		domain, err := domain.NewGLPData(
			header.Identifier,
			float64(payload.Temperature),
			float64(payload.Humidity),
			float64(payload.SoilHumidity),
			float64(payload.Light),
		)
		if err != nil {
			return nil, err
		}

		go p.db.Create(&model.GLPDataModel{
			Id: domain.Id,
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

		response, update := domain.Apply(parameters)

		if update {
			parameters.TurnOnVentilation = response.Ventilation
			parameters.TurnOnIrrigation = response.Irrigation
			parameters.TurnOnLighting = response.Lighting
			
			go p.db.Save(&parameters)
		}

		
		return response, nil
	default:
		utils.Logger.Warn().Msg("Unknown route")
		return nil, errors.New("unknown route")
	}
}

func (p *Processor) handleGet(header protocol.GLPHeader) (any, error) {
	utils.Logger.Info().Msg("Handling GET request")

	switch header.Route {
	case "parameters":
		utils.Logger.Info().Msg("Processing parameters")

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
		utils.Logger.Warn().Msg("Unknown route")
		return nil, errors.New("unknown route")
	}
}