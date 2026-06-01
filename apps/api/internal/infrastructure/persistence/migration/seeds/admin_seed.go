package seeds

import (
	"gorm.io/gorm"
)

// AdminSeeder seeds the default admin user.
type AdminSeeder struct{}

func (s *AdminSeeder) Name() string {
	return "AdminSeeder"
}

func (s *AdminSeeder) Seed(db *gorm.DB) error {
	type User struct {
		ID           string `gorm:"type:char(36);primaryKey"`
		Email        string `gorm:"type:varchar(255)"`
		PasswordHash string `gorm:"type:varchar(255)"`
		FullName     string `gorm:"type:varchar(100)"`
		Role         string `gorm:"type:enum('renter','owner','admin')"`
		IsActive     bool
	}

	admin := User{
		ID:           newUUID(),
		Email:        "admin@rentoutdoor.id",
		PasswordHash: "$2a$12$placeholder", // Replace with bcrypt hash at deploy time
		FullName:     "System Admin",
		Role:         "admin",
		IsActive:     true,
	}

	// Only seed if no admin exists
	var count int64
	db.Table("users").Where("role = ? AND deleted_at IS NULL", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	return db.Table("users").Create(&admin).Error
}
