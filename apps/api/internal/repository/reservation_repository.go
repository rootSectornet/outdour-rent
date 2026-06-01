package repository

import (
	"context"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

// ReservationRepository defines the interface for inventory reservation data access.
type ReservationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, reservation *entity.InventoryReservation) error
	FindByID(ctx context.Context, id string) (*entity.InventoryReservation, error)
	Update(ctx context.Context, reservation *entity.InventoryReservation) error
	UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status entity.ReservationStatus, updatedBy string) error

	// Availability engine
	GetOverlappingReservations(ctx context.Context, tx *gorm.DB, equipmentID string, startDate, endDate string) ([]entity.InventoryReservation, error)
	GetPeakUsage(ctx context.Context, tx *gorm.DB, equipmentID string, startDate, endDate string) (uint, error)

	// Date locks
	CreateDateLocks(ctx context.Context, tx *gorm.DB, locks []entity.ReservationDateLock) error
	DeleteDateLocksByReservation(ctx context.Context, tx *gorm.DB, reservationID string) error

	// Expiration
	FindExpired(ctx context.Context, before time.Time, limit int) ([]entity.InventoryReservation, error)
	ExpireReservation(ctx context.Context, tx *gorm.DB, id string) error
}

// OrderRepository defines the interface for order data access.
type OrderRepository interface {
	Create(ctx context.Context, tx *gorm.DB, order *entity.Order) error
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status entity.OrderStatus, updatedBy string) error
	ListByRenter(ctx context.Context, renterID string, params *OrderListParams) ([]entity.Order, *pagination.Meta, error)
	ListByStore(ctx context.Context, storeID string, params *OrderListParams) ([]entity.Order, *pagination.Meta, error)
	CreateStatusHistory(ctx context.Context, tx *gorm.DB, history *entity.OrderStatusHistory) error
}

// OrderListParams defines filtering parameters for order listing.
type OrderListParams struct {
	Status string
	pagination.Params
}

// OrderItemRepository defines the interface for order item data access.
type OrderItemRepository interface {
	Create(ctx context.Context, tx *gorm.DB, item *entity.OrderItem) error
	FindByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error)
}

// PaymentRepository defines the interface for payment data access.
type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	FindByOrderID(ctx context.Context, orderID string) ([]entity.Payment, error)
}

// DepositRepository defines the interface for deposit data access.
type DepositRepository interface {
	Create(ctx context.Context, deposit *entity.Deposit) error
	FindByOrderID(ctx context.Context, orderID string) (*entity.Deposit, error)
	Update(ctx context.Context, deposit *entity.Deposit) error
}

// RefundRepository defines the interface for refund data access.
type RefundRepository interface {
	Create(ctx context.Context, refund *entity.Refund) error
	FindByPaymentID(ctx context.Context, paymentID string) ([]entity.Refund, error)
}

// ReviewRepository defines the interface for review data access.
type ReviewRepository interface {
	Create(ctx context.Context, review *entity.Review) error
	FindByID(ctx context.Context, id string) (*entity.Review, error)
	Update(ctx context.Context, review *entity.Review) error
	ListByEquipment(ctx context.Context, equipmentID string, params *pagination.Params) ([]entity.Review, *pagination.Meta, error)
	ListByStore(ctx context.Context, storeID string, params *pagination.Params) ([]entity.Review, *pagination.Meta, error)
	ExistsByOrderAndEquipment(ctx context.Context, orderID, equipmentID string) (bool, error)
}

// NotificationRepository defines the interface for notification data access.
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	CreateBatch(ctx context.Context, notifications []entity.Notification) error
	FindByUserID(ctx context.Context, userID string, params *pagination.Params) ([]entity.Notification, *pagination.Meta, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
}

// MaintenanceRepository defines the interface for equipment maintenance data access.
type MaintenanceRepository interface {
	Create(ctx context.Context, maintenance *entity.EquipmentMaintenance) error
	FindByEquipmentID(ctx context.Context, equipmentID string) ([]entity.EquipmentMaintenance, error)
	Delete(ctx context.Context, id string) error
	GetOverlapping(ctx context.Context, tx *gorm.DB, equipmentID string, startDate, endDate string) ([]entity.EquipmentMaintenance, error)
}
