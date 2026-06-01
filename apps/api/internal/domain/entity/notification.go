package entity

import "time"

type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "in_app"
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelPush  NotificationChannel = "push"
)

type Notification struct {
	BaseModelCreateOnly
	UserID  string              `gorm:"type:char(36);not null;index:idx_notifications_user_read,priority:1;index:idx_notifications_user_type,priority:1" json:"user_id"`
	Type    string              `gorm:"type:varchar(50);not null;index:idx_notifications_user_type,priority:2" json:"type"`
	Title   string              `gorm:"type:varchar(200);not null" json:"title"`
	Body    string              `gorm:"type:text;not null" json:"body"`
	Data    *string             `gorm:"type:json" json:"data,omitempty"`
	Channel NotificationChannel `gorm:"type:enum('in_app','email','push');not null;default:'in_app'" json:"channel"`
	ReadAt  *time.Time          `gorm:"index:idx_notifications_user_read,priority:2" json:"read_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
