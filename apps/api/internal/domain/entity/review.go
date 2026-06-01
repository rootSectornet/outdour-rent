package entity

import "time"

type Review struct {
	BaseModel
	UserID         string     `gorm:"type:char(36);not null;index:idx_reviews_user_id" json:"user_id"`
	EquipmentID    string     `gorm:"type:char(36);not null;index:idx_reviews_equipment_id;uniqueIndex:uidx_reviews_order_equipment,priority:2" json:"equipment_id"`
	OrderID        string     `gorm:"type:char(36);not null;uniqueIndex:uidx_reviews_order_equipment,priority:1" json:"order_id"`
	StoreID        string     `gorm:"type:char(36);not null;index:idx_reviews_store_id" json:"store_id"`
	Rating         uint8      `gorm:"not null" json:"rating"`
	Comment        *string    `gorm:"type:text" json:"comment,omitempty"`
	OwnerReply     *string    `gorm:"type:text" json:"owner_reply,omitempty"`
	OwnerRepliedAt *time.Time `json:"owner_replied_at,omitempty"`

	// Relations
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT" json:"user,omitempty"`
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:RESTRICT" json:"equipment,omitempty"`
	Order     Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:RESTRICT" json:"order,omitempty"`
	Store     Store     `gorm:"foreignKey:StoreID;constraint:OnDelete:RESTRICT" json:"store,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}
