package models

import "github.com/google/uuid"

type Teacher struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid"`
	UserID   uuid.UUID `gorm:"type:uuid"`

	// optional profile fields
	EmployeeNo     string
	Specialization string
}