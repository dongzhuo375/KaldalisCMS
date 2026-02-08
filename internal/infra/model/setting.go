package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting stores single-row global settings.
// We use ID=1 as the singleton row.
type SystemSetting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Installed   bool       `gorm:"not null;default:false" json:"installed"`
	SiteName    string     `gorm:"type:varchar(100)" json:"site_name"`
	InstalledAt *time.Time `json:"installed_at"`
}
