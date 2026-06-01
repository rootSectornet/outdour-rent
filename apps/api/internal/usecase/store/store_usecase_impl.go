package store

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

type storeUsecase struct {
	storeRepo    repository.StoreRepository
	photoRepo    repository.StorePhotoRepository
	hoursRepo    repository.StoreOperatingHourRepository
	userRepo     repository.UserRepository
}

// NewStoreUsecase creates a new store usecase.
func NewStoreUsecase(
	storeRepo repository.StoreRepository,
	photoRepo repository.StorePhotoRepository,
	hoursRepo repository.StoreOperatingHourRepository,
	userRepo repository.UserRepository,
) usecase.StoreUsecase {
	return &storeUsecase{
		storeRepo: storeRepo,
		photoRepo: photoRepo,
		hoursRepo: hoursRepo,
		userRepo:  userRepo,
	}
}

func (uc *storeUsecase) Create(ctx context.Context, input *usecase.CreateStoreInput) (*entity.Store, error) {
	// Check if owner already has a store
	existing, err := uc.storeRepo.FindByOwnerID(ctx, input.OwnerID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: owner already has a store", usecase.ErrConflict)
	}

	// Verify user is owner role
	user, err := uc.userRepo.FindByID(ctx, input.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("%w: user not found", usecase.ErrNotFound)
	}
	if user.Role != entity.UserRoleOwner && user.Role != entity.UserRoleAdmin {
		return nil, fmt.Errorf("%w: only owners can create stores", usecase.ErrForbidden)
	}

	// Generate slug
	slug := generateSlug(input.Name)
	exists, _ := uc.storeRepo.SlugExists(ctx, slug)
	if exists {
		slug = slug + "-" + fmt.Sprintf("%d", time.Now().UnixMilli()%10000)
	}

	store := &entity.Store{
		OwnerID:     input.OwnerID,
		Name:        input.Name,
		Slug:        slug,
		Phone:       input.Phone,
		Email:       input.Email,
		Address:     input.Address,
		City:        input.City,
		Province:    input.Province,
		Status:      entity.StoreStatusPendingApproval,
	}

	if input.Description != "" {
		store.Description = &input.Description
	}
	if input.PostalCode != "" {
		store.PostalCode = &input.PostalCode
	}
	if input.Latitude != 0 {
		store.Latitude = &input.Latitude
	}
	if input.Longitude != 0 {
		store.Longitude = &input.Longitude
	}

	store.CreatedBy = &input.OwnerID

	if err := uc.storeRepo.Create(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return store, nil
}

func (uc *storeUsecase) GetByID(ctx context.Context, id string) (*entity.Store, error) {
	store, err := uc.storeRepo.FindByIDWithRelations(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: store not found", usecase.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	return store, nil
}

func (uc *storeUsecase) GetBySlug(ctx context.Context, slug string) (*entity.Store, error) {
	store, err := uc.storeRepo.FindBySlug(ctx, slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: store not found", usecase.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	return store, nil
}

func (uc *storeUsecase) GetMyStore(ctx context.Context, ownerID string) (*entity.Store, error) {
	store, err := uc.storeRepo.FindByOwnerID(ctx, ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: you don't have a store yet", usecase.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	return store, nil
}

func (uc *storeUsecase) Update(ctx context.Context, id string, input *usecase.UpdateStoreInput, ownerID string) (*entity.Store, error) {
	store, err := uc.storeRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: store not found", usecase.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find store: %w", err)
	}

	if store.OwnerID != ownerID {
		return nil, fmt.Errorf("%w: you don't own this store", usecase.ErrForbidden)
	}

	if input.Name != nil {
		store.Name = *input.Name
	}
	if input.Description != nil {
		store.Description = input.Description
	}
	if input.Phone != nil {
		store.Phone = *input.Phone
	}
	if input.Email != nil {
		store.Email = *input.Email
	}
	if input.Address != nil {
		store.Address = *input.Address
	}
	if input.City != nil {
		store.City = *input.City
	}
	if input.Province != nil {
		store.Province = *input.Province
	}
	if input.PostalCode != nil {
		store.PostalCode = input.PostalCode
	}
	if input.Latitude != nil {
		store.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		store.Longitude = input.Longitude
	}
	if input.LogoURL != nil {
		store.LogoURL = input.LogoURL
	}
	if input.BannerURL != nil {
		store.BannerURL = input.BannerURL
	}

	store.UpdatedBy = &ownerID

	if err := uc.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	return store, nil
}

func (uc *storeUsecase) List(ctx context.Context, input *usecase.ListStoreInput) ([]entity.Store, *pagination.Meta, error) {
	params := &repository.StoreListParams{
		City:     input.City,
		Province: input.Province,
		Status:   input.Status,
		Search:   input.Search,
		Params:   pagination.NewParams(input.Page, input.PerPage),
	}

	stores, meta, err := uc.storeRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list stores: %w", err)
	}

	return stores, meta, nil
}

// --- Photos ---

func (uc *storeUsecase) AddPhoto(ctx context.Context, input *usecase.AddStorePhotoInput) (*entity.StorePhoto, error) {
	store, err := uc.storeRepo.FindByID(ctx, input.StoreID)
	if err != nil {
		return nil, fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}
	if store.OwnerID != input.OwnerID {
		return nil, fmt.Errorf("%w: you don't own this store", usecase.ErrForbidden)
	}

	photo := &entity.StorePhoto{
		StoreID:   input.StoreID,
		PhotoURL:  input.PhotoURL,
		SortOrder: input.SortOrder,
		IsPrimary: input.IsPrimary,
	}
	if input.Caption != "" {
		photo.Caption = &input.Caption
	}
	photo.CreatedBy = &input.OwnerID

	// If this is primary, unset other primaries
	if input.IsPrimary {
		_ = uc.photoRepo.SetPrimary(ctx, input.StoreID, "")
	}

	if err := uc.photoRepo.Create(ctx, photo); err != nil {
		return nil, fmt.Errorf("failed to add photo: %w", err)
	}

	return photo, nil
}

func (uc *storeUsecase) RemovePhoto(ctx context.Context, storeID string, photoID string, ownerID string) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}
	if store.OwnerID != ownerID {
		return fmt.Errorf("%w: you don't own this store", usecase.ErrForbidden)
	}

	if err := uc.photoRepo.Delete(ctx, photoID); err != nil {
		return fmt.Errorf("failed to remove photo: %w", err)
	}
	return nil
}

func (uc *storeUsecase) SetPrimaryPhoto(ctx context.Context, storeID string, photoID string, ownerID string) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}
	if store.OwnerID != ownerID {
		return fmt.Errorf("%w: you don't own this store", usecase.ErrForbidden)
	}

	if err := uc.photoRepo.SetPrimary(ctx, storeID, photoID); err != nil {
		return fmt.Errorf("failed to set primary photo: %w", err)
	}
	return nil
}

// --- Operating Hours ---

func (uc *storeUsecase) SetOperatingHours(ctx context.Context, storeID string, ownerID string, hours []usecase.OperatingHourInput) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}
	if store.OwnerID != ownerID {
		return fmt.Errorf("%w: you don't own this store", usecase.ErrForbidden)
	}

	// Validate: no duplicate days
	seen := make(map[uint8]bool)
	for _, h := range hours {
		if h.DayOfWeek < 1 || h.DayOfWeek > 7 {
			return fmt.Errorf("%w: day_of_week must be 1-7", usecase.ErrValidation)
		}
		if seen[h.DayOfWeek] {
			return fmt.Errorf("%w: duplicate day_of_week %d", usecase.ErrValidation, h.DayOfWeek)
		}
		seen[h.DayOfWeek] = true
	}

	entities := make([]entity.StoreOperatingHour, len(hours))
	for i, h := range hours {
		entities[i] = entity.StoreOperatingHour{
			StoreID:   storeID,
			DayOfWeek: entity.DayOfWeek(h.DayOfWeek),
			OpenTime:  h.OpenTime,
			CloseTime: h.CloseTime,
			IsClosed:  h.IsClosed,
		}
		entities[i].CreatedBy = &ownerID
	}

	// Replace all existing hours
	_ = uc.hoursRepo.DeleteByStoreID(ctx, storeID)
	if err := uc.hoursRepo.Upsert(ctx, entities); err != nil {
		return fmt.Errorf("failed to set operating hours: %w", err)
	}
	return nil
}

func (uc *storeUsecase) GetOperatingHours(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error) {
	hours, err := uc.hoursRepo.FindByStoreID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get operating hours: %w", err)
	}
	return hours, nil
}

// --- Admin Actions ---

func (uc *storeUsecase) ApproveStore(ctx context.Context, storeID string, adminID string) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}

	if !store.IsValidStatusTransition(entity.StoreStatusActive) {
		return fmt.Errorf("%w: cannot approve store with status %s", usecase.ErrInvalidStatus, store.Status)
	}

	now := time.Now()
	store.Status = entity.StoreStatusActive
	store.VerifiedAt = &now
	store.UpdatedBy = &adminID

	if err := uc.storeRepo.Update(ctx, store); err != nil {
		return fmt.Errorf("failed to approve store: %w", err)
	}
	return nil
}

func (uc *storeUsecase) SuspendStore(ctx context.Context, storeID string, reason string, adminID string) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}

	if !store.IsValidStatusTransition(entity.StoreStatusSuspended) {
		return fmt.Errorf("%w: cannot suspend store with status %s", usecase.ErrInvalidStatus, store.Status)
	}

	now := time.Now()
	store.Status = entity.StoreStatusSuspended
	store.SuspendedAt = &now
	store.UpdatedBy = &adminID

	if err := uc.storeRepo.Update(ctx, store); err != nil {
		return fmt.Errorf("failed to suspend store: %w", err)
	}
	return nil
}

func (uc *storeUsecase) ReactivateStore(ctx context.Context, storeID string, adminID string) error {
	store, err := uc.storeRepo.FindByID(ctx, storeID)
	if err != nil {
		return fmt.Errorf("%w: store not found", usecase.ErrNotFound)
	}

	if !store.IsValidStatusTransition(entity.StoreStatusActive) {
		return fmt.Errorf("%w: cannot reactivate store with status %s", usecase.ErrInvalidStatus, store.Status)
	}

	now := time.Now()
	store.Status = entity.StoreStatusActive
	store.VerifiedAt = &now
	store.SuspendedAt = nil
	store.UpdatedBy = &adminID

	if err := uc.storeRepo.Update(ctx, store); err != nil {
		return fmt.Errorf("failed to reactivate store: %w", err)
	}
	return nil
}

// --- Helpers ---

var nonAlphaRegex = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = nonAlphaRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
