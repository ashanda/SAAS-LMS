package models

import "github.com/google/uuid"

type TenantUserRole struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantUserID uuid.UUID `gorm:"type:uuid"`
	RoleID       uuid.UUID `gorm:"type:uuid"`
}