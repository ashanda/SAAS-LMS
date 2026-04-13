package models

import "github.com/google/uuid"

type SubjectTeacher struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid"`
	SubjectID uuid.UUID `gorm:"type:uuid"`
	TeacherID uuid.UUID `gorm:"type:uuid"`
}