package models

import "github.com/google/uuid"

type Enrollment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid"`
	StudentID uuid.UUID `gorm:"type:uuid"`
	SubjectID uuid.UUID `gorm:"type:uuid"`
}