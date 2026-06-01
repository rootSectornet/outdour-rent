package mysql

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

// --- Extend StoreRepository ---

func (r *storeRepository) FindByIDWithRelations(ctx context.Context, id string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("OperatingHours", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_of_week ASC")
		}).
		Where("id = ?", id).
		First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Store{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// Update List to handle the new params
func (r *storeRepository) List(ctx context.Context, params *repository.StoreListParams) ([]entity.Store, *pagination.Meta, error) {
	var stores []entity.Store
	query := r.db.WithContext(ctx).Model(&entity.Store{})

	if params.City != "" {
		query = query.Where("city = ?", params.City)
	}
	if params.Province != "" {
		query = query.Where("province = ?", params.Province)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name LIKE ? OR city LIKE ?", search, search)
	}
	if params.OwnerID != "" {
		query = query.Where("owner_id = ?", params.OwnerID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	err := query.
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true).Limit(1)
		}).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&stores).Error
	if err != nil {
		return nil, nil, err
	}

	meta := pagination.NewMeta(total, params.Page, params.PerPage)
	return stores, meta, nil
}

// --- StorePhoto Repository ---

type storePhotoRepository struct {
	db *gorm.DB
}

// NewStorePhotoRepository creates a new store photo repository.
func NewStorePhotoRepository(db *gorm.DB) repository.StorePhotoRepository {
	return &storePhotoRepository{db: db}
}

func (r *storePhotoRepository) Create(ctx context.Context, photo *entity.StorePhoto) error {
	return r.db.WithContext(ctx).Create(photo).Error
}

func (r *storePhotoRepository) FindByStoreID(ctx context.Context, storeID string) ([]entity.StorePhoto, error) {
	var photos []entity.StorePhoto
	err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Order("sort_order ASC").
		Find(&photos).Error
	return photos, err
}

func (r *storePhotoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.StorePhoto{}, "id = ?", id).Error
}

func (r *storePhotoRepository) SetPrimary(ctx context.Context, storeID string, photoID string) error {
	// Unset all primaries for this store
	if err := r.db.WithContext(ctx).
		Model(&entity.StorePhoto{}).
		Where("store_id = ?", storeID).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	// Set the target photo as primary
	if photoID != "" {
		return r.db.WithContext(ctx).
			Model(&entity.StorePhoto{}).
			Where("id = ? AND store_id = ?", photoID, storeID).
			Update("is_primary", true).Error
	}
	return nil
}

func (r *storePhotoRepository) UpdateSortOrder(ctx context.Context, id string, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&entity.StorePhoto{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

var _ repository.StorePhotoRepository = (*storePhotoRepository)(nil)

// --- StoreOperatingHour Repository ---

type storeOperatingHourRepository struct {
	db *gorm.DB
}

// NewStoreOperatingHourRepository creates a new operating hour repository.
func NewStoreOperatingHourRepository(db *gorm.DB) repository.StoreOperatingHourRepository {
	return &storeOperatingHourRepository{db: db}
}

func (r *storeOperatingHourRepository) Upsert(ctx context.Context, hours []entity.StoreOperatingHour) error {
	if len(hours) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&hours).Error
}

func (r *storeOperatingHourRepository) FindByStoreID(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error) {
	var hours []entity.StoreOperatingHour
	err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Order("day_of_week ASC").
		Find(&hours).Error
	return hours, err
}

func (r *storeOperatingHourRepository) DeleteByStoreID(ctx context.Context, storeID string) error {
	return r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Delete(&entity.StoreOperatingHour{}).Error
}

var _ repository.StoreOperatingHourRepository = (*storeOperatingHourRepository)(nil)
