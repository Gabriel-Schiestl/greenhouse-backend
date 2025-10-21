package processor

import (
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/protocol"
	"gorm.io/gorm"
)

type Processor struct {
	db *gorm.DB
}

func NewProcessor(db *gorm.DB) *Processor {
	return &Processor{db: db}
}

func (p *Processor) Start(header protocol.GLPHeader, payload protocol.GLPPayload) map[string]any {
	switch header.Method {
		case protocol.GLPMethodPost:
			p.handlePost(header, payload)
	}

	return nil
}

func (p *Processor) handlePost(header protocol.GLPHeader, payload protocol.GLPPayload) {
	switch header.Route {
	case "sensor/data":
		
	}
}