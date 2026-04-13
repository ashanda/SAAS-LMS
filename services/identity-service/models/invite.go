package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  uuid.UUID
	Email     string
	Role      string
	Token     string
	ExpiresAt time.Time
}