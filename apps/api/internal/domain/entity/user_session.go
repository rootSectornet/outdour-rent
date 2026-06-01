package entity

import "time"

type UserSession struct {
	BaseModelCreateOnly
	UserID           string     `gorm:"type:char(36);not null;index:idx_sessions_user_id" json:"user_id"`
	RefreshTokenHash string     `gorm:"type:varchar(255);not null" json:"-"`
	UserAgent        *string    `gorm:"type:varchar(500)" json:"user_agent,omitempty"`
	IPAddress        *string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	ExpiresAt        time.Time  `gorm:"not null;index:idx_sessions_expires_at" json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

type PasswordReset struct {
	BaseModelCreateOnly
	UserID    string     `gorm:"type:char(36);not null;index:idx_password_resets_user_id" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
