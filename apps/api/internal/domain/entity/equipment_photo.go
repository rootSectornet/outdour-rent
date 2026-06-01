package entity

type EquipmentPhoto struct {
	BaseModelCreateOnly
	EquipmentID  string `gorm:"type:char(36);not null;index:idx_photos_equipment_sort,priority:1" json:"equipment_id"`
	PhotoURL     string `gorm:"type:varchar(500);not null" json:"photo_url"`
	ThumbnailURL *string `gorm:"type:varchar(500)" json:"thumbnail_url,omitempty"`
	SortOrder    int    `gorm:"not null;default:0;index:idx_photos_equipment_sort,priority:2" json:"sort_order"`
	IsPrimary    bool   `gorm:"not null;default:false" json:"is_primary"`

	// Relations
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
}

func (EquipmentPhoto) TableName() string {
	return "equipment_photos"
}

type PricingType string

const (
	PricingTypeDaily   PricingType = "daily"
	PricingTypeWeekly  PricingType = "weekly"
	PricingTypeMonthly PricingType = "monthly"
	PricingTypeCustom  PricingType = "custom"
)

type EquipmentPricing struct {
	BaseModelNoSoftDelete
	EquipmentID string      `gorm:"type:char(36);not null;index:idx_pricing_equipment_active,priority:1;index:idx_pricing_type,priority:1" json:"equipment_id"`
	PricingType PricingType `gorm:"type:enum('daily','weekly','monthly','custom');not null;index:idx_pricing_type,priority:2" json:"pricing_type"`
	MinDays     uint        `gorm:"not null;default:1" json:"min_days"`
	MaxDays     *uint       `json:"max_days,omitempty"`
	PricePerDay float64     `gorm:"type:decimal(12,2);not null" json:"price_per_day"`
	IsActive    bool        `gorm:"not null;default:true;index:idx_pricing_equipment_active,priority:2" json:"is_active"`

	// Relations
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
}

func (EquipmentPricing) TableName() string {
	return "equipment_pricing"
}

type EquipmentMaintenance struct {
	BaseModelNoSoftDelete
	EquipmentID string `gorm:"type:char(36);not null;index:idx_maintenance_equipment_dates,priority:1" json:"equipment_id"`
	Quantity    uint   `gorm:"not null" json:"quantity"`
	StartDate   string `gorm:"type:date;not null;index:idx_maintenance_equipment_dates,priority:2;index:idx_maintenance_daterange,priority:1" json:"start_date"`
	EndDate     string `gorm:"type:date;not null;index:idx_maintenance_equipment_dates,priority:3;index:idx_maintenance_daterange,priority:2" json:"end_date"`
	Reason      *string `gorm:"type:varchar(500)" json:"reason,omitempty"`

	// Relations
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
}

func (EquipmentMaintenance) TableName() string {
	return "equipment_maintenance"
}
