package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// TxManager provides transaction management.
type TxManager struct {
	db *gorm.DB
}

// NewTxManager creates a new transaction manager.
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTransaction executes fn within a database transaction.
// If fn returns an error, the transaction is rolled back.
// If fn succeeds, the transaction is committed.
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// WithTransactionIsolation executes fn within a transaction with specified isolation level.
func (tm *TxManager) WithTransactionIsolation(ctx context.Context, level string, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Set isolation level
	if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + level).Error; err != nil {
		tx.Rollback()
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func withTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
