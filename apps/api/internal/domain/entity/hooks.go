package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BeforeCreate hook generates UUID for BaseModel if not set.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook generates UUID for BaseModelNoSoftDelete if not set.
func (b *BaseModelNoSoftDelete) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook generates UUID for BaseModelCreateOnly if not set.
func (b *BaseModelCreateOnly) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
