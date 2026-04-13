package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Slug      string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Status    string         `gorm:"size:50;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	Domains   []TenantDomain `json:"domains,omitempty"`
	Branding  *TenantBranding `json:"branding,omitempty"`
}

type TenantDomain struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Domain    string    `gorm:"size:255;uniqueIndex;not null" json:"domain"`
	IsPrimary bool      `gorm:"default:true" json:"is_primary"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type TenantBranding struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"tenant_id"`
	LogoURL      string    `gorm:"type:text" json:"logo_url"`
	PrimaryColor string    `gorm:"size:20" json:"primary_color"`
	CreatedAt    time.Time `json:"created_at"`
}