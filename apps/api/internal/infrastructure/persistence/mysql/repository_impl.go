package mysql

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": deletedBy,
		}).
		Delete(&entity.User{}, "id = ?", id).Error
}

// Ensure interface compliance at compile time.
var _ repository.UserRepository = (*userRepository)(nil)
var _ repository.StoreRepository = (*storeRepository)(nil)
var _ repository.EquipmentRepository = (*equipmentRepository)(nil)
var _ repository.CategoryRepository = (*categoryRepository)(nil)

// --- Store Repository ---

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) repository.StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *entity.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindByOwnerID(ctx context.Context, ownerID string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindBySlug(ctx context.Context, slug string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) Update(ctx context.Context, store *entity.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}

func (r *storeRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Store{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_by": deletedBy}).
		Delete(&entity.Store{}, "id = ?", id).Error
}

func (r *storeRepository) List(ctx context.Context, params *repository.StoreListParams) ([]entity.Store, *pagination.Meta, error) {
	var stores []entity.Store
	query := r.db.WithContext(ctx).Model(&entity.Store{})

	if params.City != "" {
		query = query.Where("city = ?", params.City)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	query.Count(&total)

	err := query.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order("created_at DESC").
		Find(&stores).Error

	meta := pagination.NewMeta(total, params.Page, params.PerPage)
	return stores, meta, err
}

// --- Equipment Repository ---

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) repository.EquipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) Create(ctx context.Context, equipment *entity.Equipment) error {
	return r.db.WithContext(ctx).Create(equipment).Error
}

func (r *equipmentRepository) FindByID(ctx context.Context, id string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.WithContext(ctx).
		Preload("Photos", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Pricing", "is_active = ?", true).
		Preload("Category").
		Preload("Store").
		Where("id = ?", id).
		First(&equipment).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *equipmentRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := tx.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", id).
		First(&equipment).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *equipmentRepository) Update(ctx context.Context, equipment *entity.Equipment) error {
	return r.db.WithContext(ctx).Save(equipment).Error
}

func (r *equipmentRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_by": deletedBy}).
		Delete(&entity.Equipment{}, "id = ?", id).Error
}

func (r *equipmentRepository) List(ctx context.Context, params *repository.EquipmentListParams) ([]entity.Equipment, *pagination.Meta, error) {
	var items []entity.Equipment
	query := r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Preload("Photos", "is_primary = ?", true).
		Preload("Pricing", "is_active = ? AND pricing_type = ?", true, "daily")

	if params.StoreID != "" {
		query = query.Where("store_id = ?", params.StoreID)
	}
	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	var total int64
	query.Count(&total)

	err := query.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order("created_at DESC").
		Find(&items).Error

	meta := pagination.NewMeta(total, params.Page, params.PerPage)
	return items, meta, err
}

// --- Category Repository ---

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.EquipmentCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*entity.EquipmentCategory, error) {
	var cat entity.EquipmentCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.EquipmentCategory, error) {
	var cat entity.EquipmentCategory
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.EquipmentCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).
		Model(&entity.EquipmentCategory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_by": deletedBy}).
		Delete(&entity.EquipmentCategory{}, "id = ?", id).Error
}

func (r *categoryRepository) ListRoots(ctx context.Context) ([]entity.EquipmentCategory, error) {
	var cats []entity.EquipmentCategory
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC").
		Find(&cats).Error
	return cats, err
}

func (r *categoryRepository) ListByParent(ctx context.Context, parentID string) ([]entity.EquipmentCategory, error) {
	var cats []entity.EquipmentCategory
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort_order ASC").
		Find(&cats).Error
	return cats, err
}
