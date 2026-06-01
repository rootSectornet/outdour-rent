package usecase

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/pkg/pagination"
)

// RentalUsecase defines the interface for rental/order business logic.
type RentalUsecase interface {
	CreateOrder(ctx context.Context, input *CreateOrderInput) (*entity.Order, error)
	GetOrder(ctx context.Context, orderID string, userID string) (*entity.Order, error)
	ApproveOrder(ctx context.Context, orderID string, approvedBy string) error
	RejectOrder(ctx context.Context, orderID string, reason string, rejectedBy string) error
	CancelOrder(ctx context.Context, orderID string, cancelledBy string) error
	CompleteOrder(ctx context.Context, orderID string, completedBy string) error
	ActivateOrder(ctx context.Context, orderID string, activatedBy string) error
	ListRenterOrders(ctx context.Context, renterID string, params *OrderListInput) ([]entity.Order, *pagination.Meta, error)
	ListStoreOrders(ctx context.Context, storeID string, params *OrderListInput) ([]entity.Order, *pagination.Meta, error)
}

type CreateOrderInput struct {
	RenterID       string
	StoreID        string
	Items          []OrderItemInput
	RentalStartDate string
	RentalEndDate  string
	Notes          string
	PickupMethod   entity.PickupMethod
	DeliveryAddress string
}

type OrderItemInput struct {
	EquipmentID string
	Quantity    uint
}

type OrderListInput struct {
	Status string
	pagination.Params
}
