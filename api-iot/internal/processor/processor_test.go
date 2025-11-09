package processor

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Gabriel-Schiestl/greenhouse-backend/internal/model"
)

func initProcessor() *Processor {
	mockDB, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
	Conn:       mockDB,
	DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	db.AutoMigrate(model.GLPDataModel{}, model.GLPParametersModel{})
	

	processor := NewProcessor(db)
	return processor
}

// func TestHandlePost(t *testing.T) {
// 	processor := initProcessor()

// 	header := protocol.GLPHeader{
// 		PayloadLen: 34,
// 		Identifier: "DEV12345",
// 		Method:     protocol.GLPMethodPost,
// 		Route:      "sensor/data",
// 	}

// 	payload := protocol.GLPPayload{
// 		Temperature: 23.5,
// 		Humidity:    60.0,
// 		SoilHumidity: 45.0,
// 		Light:       100.0,
// 	}

// 	response, err := processor.handlePost(header, payload)

// 	require.NoError(t, err)

// }