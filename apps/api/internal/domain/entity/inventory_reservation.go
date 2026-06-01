package entity

import "time"

type ReservationStatus string

const (
	ReservationStatusPendingPayment ReservationStatus = "pending_payment"
	ReservationStatusConfirmed      ReservationStatus = "confirmed"
	ReservationStatusActive         ReservationStatus = "active"
	ReservationStatusCompleted      ReservationStatus = "completed"
	ReservationStatusCancelled      ReservationStatus = "cancelled"
	ReservationStatusExpired        ReservationStatus = "expired"
)

type InventoryReservation struct {
	BaseModelNoSoftDelete
	EquipmentID        string            `gorm:"type:char(36);not null;index:idx_reservations_equip_status_dates,priority:1;index:idx_reservations_equipment_dates,priority:1" json:"equipment_id"`
	OrderItemID        *string           `gorm:"type:char(36);uniqueIndex:idx_reservations_order_item" json:"order_item_id,omitempty"`
	Quantity           uint              `gorm:"not null" json:"quantity"`
	StartDate          string            `gorm:"type:date;not null;index:idx_reservations_equip_status_dates,priority:3;index:idx_reservations_equipment_dates,priority:2" json:"start_date"`
	EndDate            string            `gorm:"type:date;not null;index:idx_reservations_equip_status_dates,priority:4;index:idx_reservations_equipment_dates,priority:3" json:"end_date"`
	Status             ReservationStatus `gorm:"type:enum('pending_payment','confirmed','active','completed','cancelled','expired');not null;index:idx_reservations_equip_status_dates,priority:2;index:idx_reservations_status_expires,priority:1" json:"status"`
	ExpiresAt          *time.Time        `gorm:"index:idx_reservations_status_expires,priority:2" json:"expires_at,omitempty"`
	ConfirmedAt        *time.Time        `json:"confirmed_at,omitempty"`
	ActivatedAt        *time.Time        `json:"activated_at,omitempty"`
	CompletedAt        *time.Time        `json:"completed_at,omitempty"`
	CancelledAt        *time.Time        `json:"cancelled_at,omitempty"`
	CancellationReason *string           `gorm:"type:varchar(500)" json:"cancellation_reason,omitempty"`

	// Relations
	Equipment Equipment  `gorm:"foreignKey:EquipmentID;constraint:OnDelete:RESTRICT" json:"equipment,omitempty"`
	OrderItem *OrderItem `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
}

func (InventoryReservation) TableName() string {
	return "inventory_reservations"
}

type ReservationDateLock struct {
	ID            uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	EquipmentID   string  `gorm:"type:char(36);not null;uniqueIndex:uidx_date_locks_equip_date_res_maint,priority:1;index:idx_date_locks_availability,priority:1" json:"equipment_id"`
	LockDate      string  `gorm:"type:date;not null;uniqueIndex:uidx_date_locks_equip_date_res_maint,priority:2;index:idx_date_locks_availability,priority:2" json:"lock_date"`
	ReservedQty   uint    `gorm:"not null;default:0" json:"reserved_qty"`
	MaintenanceQty uint   `gorm:"not null;default:0" json:"maintenance_qty"`
	ReservationID *string `gorm:"type:char(36);uniqueIndex:uidx_date_locks_equip_date_res_maint,priority:3" json:"reservation_id,omitempty"`
	MaintenanceID *string `gorm:"type:char(36);uniqueIndex:uidx_date_locks_equip_date_res_maint,priority:4" json:"maintenance_id,omitempty"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime" json:"created_at"`

	// Relations
	Equipment   Equipment             `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
	Reservation *InventoryReservation `gorm:"foreignKey:ReservationID;constraint:OnDelete:CASCADE" json:"reservation,omitempty"`
	Maintenance *EquipmentMaintenance `gorm:"foreignKey:MaintenanceID;constraint:OnDelete:CASCADE" json:"maintenance,omitempty"`
}

func (ReservationDateLock) TableName() string {
	return "reservation_date_locks"
}
