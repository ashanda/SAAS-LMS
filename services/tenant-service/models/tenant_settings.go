package models

import "github.com/google/uuid"

type TenantSetting struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Timezone  string
	Currency  string
	Locale    string
}