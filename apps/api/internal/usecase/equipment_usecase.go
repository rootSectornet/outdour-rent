package usecase

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/pkg/pagination"
)

// EquipmentUsecase defines the interface for equipment business logic.
type EquipmentUsecase interface {
	Create(ctx context.Context, input *CreateEquipmentInput) (*entity.Equipment, error)
	GetByID(ctx context.Context, id string) (*entity.Equipment, error)
	Update(ctx context.Context, id string, input *UpdateEquipmentInput, updatedBy string) (*entity.Equipment, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	List(ctx context.Context, input *ListEquipmentInput) ([]entity.Equipment, *pagination.Meta, error)
	CheckAvailability(ctx context.Context, input *AvailabilityInput) (*AvailabilityOutput, error)
}

type CreateEquipmentInput struct {
	StoreID        string
	CategoryID     string
	Name           string
	Description    string
	Brand          string
	Specifications string
	TotalStock     uint
	Condition      entity.EquipmentCondition
	WeightGrams    uint
	MinRentalDays  uint
	MaxRentalDays  uint
	RequiresDeposit bool
	DepositAmount  float64
	CreatedBy      string
}

type UpdateEquipmentInput struct {
	Name           *string
	Description    *string
	Brand          *string
	TotalStock     *uint
	Condition      *entity.EquipmentCondition
	MinRentalDays  *uint
	MaxRentalDays  *uint
	DepositAmount  *float64
	IsActive       *bool
}

type ListEquipmentInput struct {
	StoreID    string
	CategoryID string
	City       string
	Search     string
	MinPrice   *float64
	MaxPrice   *float64
	pagination.Params
}

type AvailabilityInput struct {
	EquipmentID string
	StartDate   string
	EndDate     string
	Quantity    uint
}

type AvailabilityOutput struct {
	Available     bool
	TotalStock    uint
	PeakUsage     uint
	AvailableQty  uint
}
