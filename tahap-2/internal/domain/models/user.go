package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"`
	FirstName   string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName    string         `gorm:"type:varchar(100)" json:"last_name"`
	PhoneNumber string         `gorm:"type:varchar(20);unique;not null" json:"phone_number"`
	Address     string         `gorm:"type:text" json:"address"`
	PIN         string         `gorm:"type:varchar(255);not null" json:"-"`
	Balance     float64        `gorm:"type:decimal(15,2);default:0" json:"balance"`
	CreatedAt   time.Time      `json:"created_date"`
	UpdatedAt   time.Time      `json:"updated_date"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
