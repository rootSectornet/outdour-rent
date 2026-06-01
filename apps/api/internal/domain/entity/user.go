package entity

import "time"

type User struct {
	BaseModel
	Email           string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL" json:"email"`
	PasswordHash    string     `gorm:"type:varchar(255)" json:"-"`
	FullName        string     `gorm:"type:varchar(100);not null" json:"full_name"`
	Phone           *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL       *string    `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	Role            UserRole   `gorm:"type:enum('renter','owner','admin');not null;default:'renter';index:idx_users_role" json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	IsActive        bool       `gorm:"not null;default:true" json:"is_active"`

	// OAuth
	GoogleID *string `gorm:"type:varchar(255);uniqueIndex:idx_users_google_id" json:"-"`
	Provider string  `gorm:"type:varchar(50);not null;default:'local'" json:"provider"`

	// Relations
	Stores        []Store        `gorm:"foreignKey:OwnerID" json:"stores,omitempty"`
	Orders        []Order        `gorm:"foreignKey:RenterID" json:"orders,omitempty"`
	Reviews       []Review       `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Sessions      []UserSession  `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type UserRole string

const (
	UserRoleRenter UserRole = "renter"
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
)
