package models

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code string    `gorm:"uniqueIndex"`
	Name string
}