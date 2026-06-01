package entity

import "time"

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusApproved       OrderStatus = "approved"
	OrderStatusRejected       OrderStatus = "rejected"
	OrderStatusActive         OrderStatus = "active"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusExpired        OrderStatus = "expired"
)

type PickupMethod string

const (
	PickupMethodPickup   PickupMethod = "pickup"
	PickupMethodDelivery PickupMethod = "delivery"
)

type Order struct {
	BaseModel
	OrderNumber     string       `gorm:"type:varchar(20);not null;uniqueIndex:idx_orders_order_number" json:"order_number"`
	RenterID        string       `gorm:"type:char(36);not null;index:idx_orders_renter_status,priority:1" json:"renter_id"`
	StoreID         string       `gorm:"type:char(36);not null;index:idx_orders_store_status,priority:1" json:"store_id"`
	Status          OrderStatus  `gorm:"type:enum('pending_payment','paid','approved','rejected','active','completed','cancelled','expired');not null;index:idx_orders_renter_status,priority:2;index:idx_orders_store_status,priority:2;index:idx_orders_status_deadline,priority:1" json:"status"`
	RentalStartDate string       `gorm:"type:date;not null;index:idx_orders_dates,priority:1" json:"rental_start_date"`
	RentalEndDate   string       `gorm:"type:date;not null;index:idx_orders_dates,priority:2" json:"rental_end_date"`
	RentalDays      uint         `gorm:"not null" json:"rental_days"`
	Subtotal        float64      `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	ServiceFee      float64      `gorm:"type:decimal(12,2);not null;default:0.00" json:"service_fee"`
	DepositTotal    float64      `gorm:"type:decimal(12,2);not null;default:0.00" json:"deposit_total"`
	TotalAmount     float64      `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Notes           *string      `gorm:"type:text" json:"notes,omitempty"`
	PickupMethod    PickupMethod `gorm:"type:enum('pickup','delivery');not null;default:'pickup'" json:"pickup_method"`
	DeliveryAddress *string      `gorm:"type:text" json:"delivery_address,omitempty"`
	PaymentDeadline time.Time    `gorm:"not null;index:idx_orders_status_deadline,priority:2" json:"payment_deadline"`
	ApprovedAt      *time.Time   `json:"approved_at,omitempty"`
	RejectedAt      *time.Time   `json:"rejected_at,omitempty"`
	RejectionReason *string      `gorm:"type:varchar(500)" json:"rejection_reason,omitempty"`
	CompletedAt     *time.Time   `json:"completed_at,omitempty"`
	CancelledAt     *time.Time   `json:"cancelled_at,omitempty"`

	// Relations
	Renter        User               `gorm:"foreignKey:RenterID;constraint:OnDelete:RESTRICT" json:"renter,omitempty"`
	Store         Store              `gorm:"foreignKey:StoreID;constraint:OnDelete:RESTRICT" json:"store,omitempty"`
	Items         []OrderItem        `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payments      []Payment          `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	Deposits      []Deposit          `gorm:"foreignKey:OrderID" json:"deposits,omitempty"`
	StatusHistory []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"status_history,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	BaseModelNoSoftDelete
	OrderID         string  `gorm:"type:char(36);not null;index:idx_order_items_order_id" json:"order_id"`
	EquipmentID     string  `gorm:"type:char(36);not null;index:idx_order_items_equipment_id" json:"equipment_id"`
	ReservationID   string  `gorm:"type:char(36);not null;uniqueIndex:idx_order_items_reservation_id" json:"reservation_id"`
	Quantity        uint    `gorm:"not null" json:"quantity"`
	PricePerDay     float64 `gorm:"type:decimal(12,2);not null" json:"price_per_day"`
	RentalDays      uint    `gorm:"not null" json:"rental_days"`
	Subtotal        float64 `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	DepositPerUnit  float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"deposit_per_unit"`
	DepositSubtotal float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"deposit_subtotal"`

	// Relations
	Order       Order                `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Equipment   Equipment            `gorm:"foreignKey:EquipmentID;constraint:OnDelete:RESTRICT" json:"equipment,omitempty"`
	Reservation InventoryReservation `gorm:"foreignKey:ReservationID;constraint:OnDelete:RESTRICT" json:"reservation,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type OrderStatusHistory struct {
	BaseModelCreateOnly
	OrderID    string  `gorm:"type:char(36);not null;index:idx_status_history_order_created,priority:1" json:"order_id"`
	FromStatus *string `gorm:"type:varchar(30)" json:"from_status,omitempty"`
	ToStatus   string  `gorm:"type:varchar(30);not null" json:"to_status"`
	ChangedBy  *string `gorm:"type:char(36)" json:"changed_by,omitempty"`
	Reason     *string `gorm:"type:varchar(500)" json:"reason,omitempty"`
	Metadata   *string `gorm:"type:json" json:"metadata,omitempty"`

	// Relations
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
}

func (OrderStatusHistory) TableName() string {
	return "order_status_history"
}
