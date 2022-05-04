package seed

import (
	"OMUS/server/models"
	"log"

	"github.com/jinzhu/gorm"
)

/*
	Create DB entity and fill it with raw test data
*/

var urls = []models.URL{
	models.URL{
		OriginalURL:   "7",
		EncodedURL:    "h",
		VisitsCounter: 0,
	},
	models.URL{
		OriginalURL:   "555555",
		EncodedURL:    "JGuc",
		VisitsCounter: 0,
	},
	models.URL{
		OriginalURL:   "10000000000",
		EncodedURL:    "KY8U4k",
		VisitsCounter: 0,
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.URL{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.URL{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range urls {
		err = db.Debug().Model(&models.URL{}).Create(&urls[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}
}
