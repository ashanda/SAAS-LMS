package models

import "github.com/google/uuid"

type Subject struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID uuid.UUID
	Name     string
	GradeID  uuid.UUID
}