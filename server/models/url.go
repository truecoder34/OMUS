package models

import (
	helper "OMUS/server/helpers"
	"errors"
	"html"
	"math/rand"
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
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

/*
	Method will be called BEFORE each create call in ORM.
 	And will generate UUID for ID
*/
func (base *Entity) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

/*
	Classic URL entity storing data about original and encoded data
*/
type URL struct {
	Entity
	OriginalURL        string    `gorm:"size:255;not null;unique" json:"originalURL"`
	EncodedURL         string    `gorm:"size:255;not null;unique" json:"encodedURL"`
	EOL                time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"eol"`
	VisitsCounter      int64     `gorm:"not null;default:0" json:"visitsCounter"`
	RegeneratesCounter int64     `gorm:"not null;default:0" json:"regeneratesCounter"`
}

/*
	Method to create empty URL entity ready to use
*/
func (url *URL) Prepare() {
	url.Entity = Entity{}
	url.OriginalURL = html.EscapeString(strings.TrimSpace(url.OriginalURL))
	RandomIntegerwithinRange := rand.Uint64()
	url.EncodedURL = strings.TrimSpace(helper.Encode(RandomIntegerwithinRange))
	url.EOL = time.Now().Add(time.Duration(7776000) * time.Second) // life time = 90 days = 90d*24h*60m*60s = 7776000 s
	url.VisitsCounter = 0
	url.RegeneratesCounter = 0
}

/*
	URL Entites Fields validator
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
	Validate URL by existence . Check BY ENCODED URL
*/
func (url *URL) ValidateOnExistence(db *gorm.DB, originalURL string) bool {
	// GET ENTITY BY ENCODED URL
	entity := URL{}
	var err error = db.Debug().Model(&URL{}).Where("original_url = ?", originalURL).Take(&entity).Error
	if err != nil {
		return false
	}
	if (entity == URL{}) {
		return false
	}
	return true

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
	Get entity by encoded URL
*/
func (url *URL) GetEntityByEncodedURL(db *gorm.DB, encodedURL string) (*URL, error) {
	var err error = db.Debug().Model(&URL{}).Where("encoded_url = ?", encodedURL).Take(&url).Error
	if err != nil {
		return &URL{}, err
	}

	return url, nil
}

/*
	Get entity by original URL
*/
func (url *URL) GetEntityByOriginalURL(db *gorm.DB, originalURL string) (*URL, error) {
	var err error = db.Debug().Model(&URL{}).Where("original_url = ?", originalURL).Take(&url).Error
	if err != nil {
		return &URL{}, err
	}

	return url, nil
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

/*
	Update required URL entity by ID : INCREMENT  VISITS COUNTER
*/
func (url *URL) UpdateURL(db *gorm.DB, mode int) (*URL, error) {

	url.OriginalURL = html.EscapeString(strings.TrimSpace(url.OriginalURL))
	url.EncodedURL = html.EscapeString(strings.TrimSpace(url.EncodedURL))
	url.EOL = time.Now().Add(time.Duration(7776000) * time.Second) // life time = 90 days = 90d*24h*60m*60s = 7776000 s

	if mode == 1 {
		// redirect mode
		url.VisitsCounter += 1
	} else if mode == 2 {
		// regenerate mode
		url.RegeneratesCounter += 1
	} else {
		url.VisitsCounter = 0
		url.RegeneratesCounter = 0
	}

	var err error = db.Debug().Model(&URL{}).Where("id = ?", url.ID).Updates(URL{
		OriginalURL:        url.OriginalURL,
		EncodedURL:         url.EncodedURL,
		EOL:                time.Now().Add(time.Duration(7776000) * time.Second),
		VisitsCounter:      url.VisitsCounter,
		RegeneratesCounter: url.RegeneratesCounter,
	}).Error
	if err != nil {
		return &URL{}, err
	}

	return url, nil
}

/*
	Update required URL entity by ID : INCREMENT  VISITS COUNTER
	Increment REGENERATIONS COUNTER
*/
func (url *URL) IncrementURLRegenerations(db *gorm.DB, encodedURL string) (*URL, error) {

	existingURL := URL{}
	var err error = db.Debug().Model(&URL{}).Where("encoded_url = ?", encodedURL).Take(&existingURL).Error
	if err != nil {
		return url, err
	}

	err = db.Debug().Model(&URL{}).Where("id = ?", existingURL.ID).Updates(URL{
		OriginalURL:        existingURL.OriginalURL,
		EncodedURL:         existingURL.EncodedURL,
		EOL:                existingURL.EOL,
		VisitsCounter:      existingURL.VisitsCounter,
		RegeneratesCounter: existingURL.RegeneratesCounter + 1,
	}).Error
	if err != nil {
		return &URL{}, err
	}

	return url, nil
}

/*
	Delete URL note by ID
*/
func (p *URL) DeleteURL(db *gorm.DB, pid uuid.UUID) (int64, error) {
	db = db.Debug().Model(&URL{}).Where("id = ?", pid).Take(&URL{}).Delete(&URL{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("URL Entity  not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
