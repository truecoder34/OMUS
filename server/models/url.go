package models

import (
	"errors"
	"html"
	"strings"
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

/*
	Method to create empty URL entity ready to use
*/
func (url *URL) Prepare() {
	url.Entity = Entity{}
	url.OriginalURL = html.EscapeString(strings.TrimSpace(url.OriginalURL))
	url.EncodedURL = html.EscapeString(strings.TrimSpace(url.EncodedURL))
	url.EOL = time.Now().Add(time.Duration(7776000) * time.Second) // life time = 90 days = 90d*24h*60m*60s = 7776000 s
	url.VisitsCounter = 0
}

/*
	UDL Entites Fields validator
*/
func (url *URL) Validate() error {
	if url.OriginalURL == "" {
		return errors.New("URL data is required;")
	}
	if url.EncodedURL == "" {
		return errors.New("encoded URL data is required;")
	}

	return nil
}

/*
	Save URL Entity to DB
*/
func (url *URL) SaveURL(db *gorm.DB) (*URL, error) {
	var err error = db.Debug().Model(&URL{}).Create(&url).Error
	if err != nil {
		return &URL{}, err
	}

	return url, nil
}

/*
	Get all URL notes from DB
*/
func (url *URL) FindAllURLs(db *gorm.DB) (*[]URL, error) {
	urls := []URL{}
	var err error = db.Debug().Model(&URL{}).Limit(100).Find(&urls).Error
	if err != nil {
		return &[]URL{}, err
	}
	return &urls, nil
}

/*
	Get get URL note from DB by ID
*/
func (url *URL) FindURLbyID(db *gorm.DB, pid uuid.UUID) (*URL, error) {
	var err error = db.Debug().Model(&URL{}).Where("id = ?", pid).Take(&url).Error
	if err != nil {
		return &URL{}, err
	}

	if url.ID != uuid.Nil {
		err = db.Debug().Model(&URL{}).Where("id = ?", url.ID).Take(&url.ID).Error
		if err != nil {
			return &URL{}, err
		}
	}
	return url, nil
}
