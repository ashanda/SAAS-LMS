package models

import "github.com/google/uuid"

type TenantUser struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid"`
	UserID   uuid.UUID `gorm:"type:uuid"` 
	Status   string
}