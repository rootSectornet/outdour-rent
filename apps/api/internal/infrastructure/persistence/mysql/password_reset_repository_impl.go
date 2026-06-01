package mysql

import (
	"context"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository instance.
func NewPasswordResetRepository(db *gorm.DB) repository.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, reset *entity.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordReset, error) {
	var reset entity.PasswordReset
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.PasswordReset{}).
		Where("id = ?", id).
		Update("used_at", &now).Error
}

func (r *passwordResetRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entity.PasswordReset{}).Error
}

// Ensure interface compliance at compile time.
var _ repository.PasswordResetRepository = (*passwordResetRepository)(nil)
