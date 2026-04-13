package models

import "github.com/google/uuid"

type Student struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID uuid.UUID
	UserID   uuid.UUID
}