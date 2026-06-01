package equipment

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

type equipmentUsecase struct {
	equipmentRepo   repository.EquipmentRepository
	categoryRepo    repository.CategoryRepository
	reservationRepo repository.ReservationRepository
	storeRepo       repository.StoreRepository
	txManager       repository.TransactionManager
}

// NewEquipmentUsecase creates a new equipment usecase.
func NewEquipmentUsecase(
	equipmentRepo repository.EquipmentRepository,
	categoryRepo repository.CategoryRepository,
	reservationRepo repository.ReservationRepository,
	storeRepo repository.StoreRepository,
	txManager repository.TransactionManager,
) usecase.EquipmentUsecase {
	return &equipmentUsecase{
		equipmentRepo:   equipmentRepo,
		categoryRepo:    categoryRepo,
		reservationRepo: reservationRepo,
		storeRepo:       storeRepo,
		txManager:       txManager,
	}
}

func (uc *equipmentUsecase) Create(ctx context.Context, input *usecase.CreateEquipmentInput) (*entity.Equipment, error) {
	// Verify store exists and user owns it
	store, err := uc.storeRepo.FindByOwnerID(ctx, input.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w: you don't have a store", usecase.ErrForbidden)
	}

	// Verify category exists
	_, err = uc.categoryRepo.FindByID(ctx, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: category not found", usecase.ErrNotFound)
	}

	// Set defaults
	condition := input.Condition
	if condition == "" {
		condition = entity.EquipmentConditionGood
	}
	minRental := input.MinRentalDays
	if minRental == 0 {
		minRental = 1
	}
	maxRental := input.MaxRentalDays
	if maxRental == 0 {
		maxRental = 30
	}

	slug := generateSlug(input.Name)

	equipment := &entity.Equipment{
		StoreID:         store.ID,
		CategoryID:      input.CategoryID,
		Name:            input.Name,
		Slug:            slug,
		Condition:       condition,
		Status:          entity.EquipmentStatusAvailable,
		TotalStock:      input.TotalStock,
		PurchaseDate:    input.PurchaseDate,
		MinRentalDays:   minRental,
		MaxRentalDays:   maxRental,
		RequiresDeposit: input.RequiresDeposit,
		DepositAmount:   input.DepositAmount,
		IsActive:        true,
	}

	if input.Description != "" {
		equipment.Description = &input.Description
	}
	if input.Brand != "" {
		equipment.Brand = &input.Brand
	}
	if input.Specifications != "" {
		equipment.Specifications = &input.Specifications
	}
	if input.WeightGrams > 0 {
		equipment.WeightGrams = &input.WeightGrams
	}

	equipment.BaseModel.CreatedBy = &input.CreatedBy

	if err := uc.equipmentRepo.Create(ctx, equipment); err != nil {
		return nil, err
	}

	return uc.equipmentRepo.FindByID(ctx, equipment.ID)
}

func (uc *equipmentUsecase) GetByID(ctx context.Context, id string) (*entity.Equipment, error) {
	equipment, err := uc.equipmentRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: equipment not found", usecase.ErrNotFound)
		}
		return nil, err
	}
	return equipment, nil
}

func (uc *equipmentUsecase) Update(ctx context.Context, id string, input *usecase.UpdateEquipmentInput, updatedBy string) (*entity.Equipment, error) {
	equipment, err := uc.equipmentRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: equipment not found", usecase.ErrNotFound)
		}
		return nil, err
	}

	// Verify ownership via store
	store, err := uc.storeRepo.FindByOwnerID(ctx, updatedBy)
	if err != nil || store.ID != equipment.StoreID {
		return nil, fmt.Errorf("%w: you don't own this equipment", usecase.ErrForbidden)
	}

	// Apply updates
	if input.Name != nil {
		equipment.Name = *input.Name
		equipment.Slug = generateSlug(*input.Name)
	}
	if input.Description != nil {
		equipment.Description = input.Description
	}
	if input.Brand != nil {
		equipment.Brand = input.Brand
	}
	if input.TotalStock != nil {
		equipment.TotalStock = *input.TotalStock
	}
	if input.Condition != nil {
		equipment.Condition = *input.Condition
	}
	if input.PurchaseDate != nil {
		equipment.PurchaseDate = input.PurchaseDate
	}
	if input.MinRentalDays != nil {
		equipment.MinRentalDays = *input.MinRentalDays
	}
	if input.MaxRentalDays != nil {
		equipment.MaxRentalDays = *input.MaxRentalDays
	}
	if input.DepositAmount != nil {
		equipment.DepositAmount = *input.DepositAmount
	}
	if input.IsActive != nil {
		equipment.IsActive = *input.IsActive
	}

	equipment.BaseModel.UpdatedBy = &updatedBy

	if err := uc.equipmentRepo.Update(ctx, equipment); err != nil {
		return nil, err
	}

	return uc.equipmentRepo.FindByID(ctx, equipment.ID)
}

func (uc *equipmentUsecase) Delete(ctx context.Context, id string, deletedBy string) error {
	equipment, err := uc.equipmentRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%w: equipment not found", usecase.ErrNotFound)
		}
		return err
	}

	// Verify ownership
	store, err := uc.storeRepo.FindByOwnerID(ctx, deletedBy)
	if err != nil || store.ID != equipment.StoreID {
		return fmt.Errorf("%w: you don't own this equipment", usecase.ErrForbidden)
	}

	return uc.equipmentRepo.Delete(ctx, id, deletedBy)
}

func (uc *equipmentUsecase) List(ctx context.Context, input *usecase.ListEquipmentInput) ([]entity.Equipment, *pagination.Meta, error) {
	isActive := true
	params := &repository.EquipmentListParams{
		StoreID:    input.StoreID,
		CategoryID: input.CategoryID,
		City:       input.City,
		Search:     input.Search,
		Status:     input.Status,
		MinPrice:   input.MinPrice,
		MaxPrice:   input.MaxPrice,
		IsActive:   &isActive,
		Params:     input.Params,
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PerPage == 0 {
		params.PerPage = 20
	}

	return uc.equipmentRepo.List(ctx, params)
}

func (uc *equipmentUsecase) ChangeStatus(ctx context.Context, id string, newStatus entity.EquipmentStatus, updatedBy string) (*entity.Equipment, error) {
	equipment, err := uc.equipmentRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: equipment not found", usecase.ErrNotFound)
		}
		return nil, err
	}

	// Verify ownership
	store, err := uc.storeRepo.FindByOwnerID(ctx, updatedBy)
	if err != nil || store.ID != equipment.StoreID {
		return nil, fmt.Errorf("%w: you don't own this equipment", usecase.ErrForbidden)
	}

	// Validate status transition
	if !entity.IsValidEquipmentStatusTransition(equipment.Status, newStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", usecase.ErrInvalidStatus, equipment.Status, newStatus)
	}

	equipment.Status = newStatus
	equipment.BaseModel.UpdatedBy = &updatedBy

	// If retired, also deactivate
	if newStatus == entity.EquipmentStatusRetired {
		equipment.IsActive = false
	}

	if err := uc.equipmentRepo.Update(ctx, equipment); err != nil {
		return nil, err
	}

	return uc.equipmentRepo.FindByID(ctx, equipment.ID)
}

func (uc *equipmentUsecase) CheckAvailability(ctx context.Context, input *usecase.AvailabilityInput) (*usecase.AvailabilityOutput, error) {
	equipment, err := uc.equipmentRepo.FindByID(ctx, input.EquipmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: equipment not found", usecase.ErrNotFound)
		}
		return nil, err
	}

	// Equipment must be active and available status to be rented
	if !equipment.IsActive || equipment.Status == entity.EquipmentStatusRetired || equipment.Status == entity.EquipmentStatusDamaged {
		return &usecase.AvailabilityOutput{
			Available:    false,
			TotalStock:   equipment.TotalStock,
			PeakUsage:    0,
			AvailableQty: 0,
			Status:       string(equipment.Status),
		}, nil
	}

	// Validate dates
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid start_date format", usecase.ErrValidation)
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid end_date format", usecase.ErrValidation)
	}
	if !endDate.After(startDate) {
		return nil, fmt.Errorf("%w: end_date must be after start_date", usecase.ErrValidation)
	}

	quantity := input.Quantity
	if quantity == 0 {
		quantity = 1
	}

	// Check peak usage across requested date range using transaction-safe method
	var peakUsage uint
	err = uc.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		var txErr error
		peakUsage, txErr = uc.reservationRepo.GetPeakUsage(ctx, tx, input.EquipmentID, input.StartDate, input.EndDate)
		return txErr
	})
	if err != nil {
		// If no reservations exist yet, peak usage is 0
		peakUsage = 0
	}

	availableQty := equipment.TotalStock - peakUsage
	if peakUsage > equipment.TotalStock {
		availableQty = 0
	}

	return &usecase.AvailabilityOutput{
		Available:    availableQty >= quantity,
		TotalStock:   equipment.TotalStock,
		PeakUsage:    peakUsage,
		AvailableQty: availableQty,
		Status:       string(equipment.Status),
	}, nil
}

// --- Helpers ---

var nonAlphaNumRegex = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = nonAlphaNumRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
