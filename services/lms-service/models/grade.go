package models

import "github.com/google/uuid"

type Grade struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid"`
	Name     string
}