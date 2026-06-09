package entity

import "time"

type PasswordResetToken struct {
	ID         uint64     `gorm:"primaryKey"`
	UserID     string     `gorm:"type:char(36);not null;index"`
	Token      string     `gorm:"size:128;not null;uniqueIndex"`
	ExpiredAt  time.Time  `gorm:"column:expired_at;not null"`
	UsedAt     *time.Time `gorm:"column:used_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
