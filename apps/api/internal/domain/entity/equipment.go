package entity

import "time"

type EquipmentCondition string

const (
	EquipmentConditionNew  EquipmentCondition = "new"
	EquipmentConditionGood EquipmentCondition = "good"
	EquipmentConditionFair EquipmentCondition = "fair"
	EquipmentConditionWorn EquipmentCondition = "worn"
)

type EquipmentStatus string

const (
	EquipmentStatusAvailable   EquipmentStatus = "available"
	EquipmentStatusReserved    EquipmentStatus = "reserved"
	EquipmentStatusRented      EquipmentStatus = "rented"
	EquipmentStatusMaintenance EquipmentStatus = "maintenance"
	EquipmentStatusDamaged     EquipmentStatus = "damaged"
	EquipmentStatusRetired     EquipmentStatus = "retired"
)

// IsValidEquipmentStatusTransition checks if a status transition is allowed.
func IsValidEquipmentStatusTransition(from, to EquipmentStatus) bool {
	transitions := map[EquipmentStatus][]EquipmentStatus{
		EquipmentStatusAvailable:   {EquipmentStatusReserved, EquipmentStatusRented, EquipmentStatusMaintenance, EquipmentStatusDamaged, EquipmentStatusRetired},
		EquipmentStatusReserved:    {EquipmentStatusAvailable, EquipmentStatusRented, EquipmentStatusMaintenance, EquipmentStatusDamaged},
		EquipmentStatusRented:      {EquipmentStatusAvailable, EquipmentStatusMaintenance, EquipmentStatusDamaged},
		EquipmentStatusMaintenance: {EquipmentStatusAvailable, EquipmentStatusDamaged, EquipmentStatusRetired},
		EquipmentStatusDamaged:     {EquipmentStatusMaintenance, EquipmentStatusRetired},
		EquipmentStatusRetired:     {},
	}

	allowed, exists := transitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

type Equipment struct {
	BaseModel
	StoreID         string             `gorm:"type:char(36);not null;index:idx_equipment_store_id;uniqueIndex:idx_equipment_store_slug" json:"store_id"`
	CategoryID      string             `gorm:"type:char(36);not null;index:idx_equipment_category_id" json:"category_id"`
	Name            string             `gorm:"type:varchar(200);not null" json:"name"`
	Slug            string             `gorm:"type:varchar(220);not null;uniqueIndex:idx_equipment_store_slug" json:"slug"`
	Description     *string            `gorm:"type:text" json:"description,omitempty"`
	Brand           *string            `gorm:"type:varchar(100)" json:"brand,omitempty"`
	Specifications  *string            `gorm:"type:json" json:"specifications,omitempty"`
	TotalStock      uint               `gorm:"not null;default:0" json:"total_stock"`
	Condition       EquipmentCondition `gorm:"type:enum('new','good','fair','worn');not null;default:'good'" json:"condition"`
	Status          EquipmentStatus    `gorm:"type:enum('available','reserved','rented','maintenance','damaged','retired');not null;default:'available';index:idx_equipment_status" json:"status"`
	PurchaseDate    *time.Time         `gorm:"type:date" json:"purchase_date,omitempty"`
	WeightGrams     *uint              `json:"weight_grams,omitempty"`
	MinRentalDays   uint               `gorm:"not null;default:1" json:"min_rental_days"`
	MaxRentalDays   uint               `gorm:"not null;default:30" json:"max_rental_days"`
	RequiresDeposit bool               `gorm:"not null;default:true" json:"requires_deposit"`
	DepositAmount   float64            `gorm:"type:decimal(12,2);not null;default:0.00" json:"deposit_amount"`
	IsActive        bool               `gorm:"not null;default:true;index:idx_equipment_active" json:"is_active"`
	RatingAvg       float64            `gorm:"type:decimal(3,2);not null;default:0.00" json:"rating_avg"`
	RatingCount     uint               `gorm:"not null;default:0" json:"rating_count"`
	RentalCount     uint               `gorm:"not null;default:0" json:"rental_count"`

	// Relations
	Store        Store                  `gorm:"foreignKey:StoreID;constraint:OnDelete:RESTRICT" json:"store,omitempty"`
	Category     EquipmentCategory      `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category,omitempty"`
	Photos       []EquipmentPhoto       `gorm:"foreignKey:EquipmentID" json:"photos,omitempty"`
	Pricing      []EquipmentPricing     `gorm:"foreignKey:EquipmentID" json:"pricing,omitempty"`
	Maintenance  []EquipmentMaintenance `gorm:"foreignKey:EquipmentID" json:"maintenance,omitempty"`
	Reservations []InventoryReservation `gorm:"foreignKey:EquipmentID" json:"reservations,omitempty"`
}

func (Equipment) TableName() string {
	return "equipment"
}
