package repository

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

// TransactionManager defines the interface for transaction management.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTransactionIsolation(ctx context.Context, level string, fn func(tx *gorm.DB) error) error
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string, deletedBy string) error
}

// SessionRepository defines the interface for session management.
type SessionRepository interface {
	Create(ctx context.Context, session *entity.UserSession) error
	FindByID(ctx context.Context, id string) (*entity.UserSession, error)
	FindActiveByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.UserSession, error)
	Revoke(ctx context.Context, id string) error
	RevokeByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}

// PasswordResetRepository defines the interface for password reset tokens.
type PasswordResetRepository interface {
	Create(ctx context.Context, reset *entity.PasswordReset) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordReset, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

// StoreRepository defines the interface for store data access.
type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	FindByID(ctx context.Context, id string) (*entity.Store, error)
	FindByIDWithRelations(ctx context.Context, id string) (*entity.Store, error)
	FindByOwnerID(ctx context.Context, ownerID string) (*entity.Store, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Store, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
	Update(ctx context.Context, store *entity.Store) error
	Delete(ctx context.Context, id string, deletedBy string) error
	List(ctx context.Context, params *StoreListParams) ([]entity.Store, *pagination.Meta, error)
}

// StoreListParams defines filtering parameters for store listing.
type StoreListParams struct {
	City     string
	Province string
	Status   string
	Search   string
	OwnerID  string
	pagination.Params
}

// StorePhotoRepository defines the interface for store photos.
type StorePhotoRepository interface {
	Create(ctx context.Context, photo *entity.StorePhoto) error
	FindByStoreID(ctx context.Context, storeID string) ([]entity.StorePhoto, error)
	Delete(ctx context.Context, id string) error
	SetPrimary(ctx context.Context, storeID string, photoID string) error
	UpdateSortOrder(ctx context.Context, id string, sortOrder int) error
}

// StoreOperatingHourRepository defines the interface for store hours.
type StoreOperatingHourRepository interface {
	Upsert(ctx context.Context, hours []entity.StoreOperatingHour) error
	FindByStoreID(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error)
	DeleteByStoreID(ctx context.Context, storeID string) error
}

// EquipmentRepository defines the interface for equipment data access.
type EquipmentRepository interface {
	Create(ctx context.Context, equipment *entity.Equipment) error
	FindByID(ctx context.Context, id string) (*entity.Equipment, error)
	Update(ctx context.Context, equipment *entity.Equipment) error
	Delete(ctx context.Context, id string, deletedBy string) error
	List(ctx context.Context, params *EquipmentListParams) ([]entity.Equipment, *pagination.Meta, error)
	FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id string) (*entity.Equipment, error)
}

// EquipmentListParams defines filtering parameters for equipment listing.
type EquipmentListParams struct {
	StoreID    string
	CategoryID string
	City       string
	Search     string
	MinPrice   *float64
	MaxPrice   *float64
	IsActive   *bool
	pagination.Params
}

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.EquipmentCategory) error
	FindByID(ctx context.Context, id string) (*entity.EquipmentCategory, error)
	FindBySlug(ctx context.Context, slug string) (*entity.EquipmentCategory, error)
	Update(ctx context.Context, category *entity.EquipmentCategory) error
	Delete(ctx context.Context, id string, deletedBy string) error
	ListRoots(ctx context.Context) ([]entity.EquipmentCategory, error)
	ListByParent(ctx context.Context, parentID string) ([]entity.EquipmentCategory, error)
}
