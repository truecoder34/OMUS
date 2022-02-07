package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

/*
	Create base entity
*/
type Entity struct {
	//ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ID        uuid.UUID  `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
}

/*
	Method will be called BEFORE each create call in ORM.
 	And wiil generate UUID for ID
*/
func (base *Entity) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

/*
	Classic URL entity storing data aboiut original and encoded data
*/
type URL struct {
	Entity
	OriginalURL   string    `gorm:"size:255;not null;unique" json:"originalURL"`
	EncodedURL    string    `gorm:"size:255;not null;unique" json:"encodedURL"`
	EOL           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"eol"`
	VisitsCounter int64     `gorm:"not null;default:0;unique" json:"visitsCounter"`
}
