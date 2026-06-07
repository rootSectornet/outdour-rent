package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel provides UUID primary key with full audit trail.
type BaseModel struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	CreatedBy *string        `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy *string        `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedBy *string        `gorm:"type:char(36)" json:"deleted_by,omitempty"`
}

// BaseModelNoSoftDelete provides UUID primary key without soft delete.
type BaseModelNoSoftDelete struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
	CreatedBy *string   `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy *string   `gorm:"type:char(36)" json:"updated_by,omitempty"`
}

// BaseModelCreateOnly provides UUID primary key with created_at only.
type BaseModelCreateOnly struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"type:char(36)" json:"created_by,omitempty"`
}
