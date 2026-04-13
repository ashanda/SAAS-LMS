package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string    `gorm:"size:255"`
	Username  string    `gorm:"size:100;uniqueIndex"`
	Email     string    `gorm:"size:255;uniqueIndex;not null"`
	Phone     string    `gorm:"size:30"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}