package mysql

import (
	"context"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"gorm.io/gorm"
)

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository instance.
func NewReservationRepository(db *gorm.DB) repository.ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, tx *gorm.DB, reservation *entity.InventoryReservation) error {
	return tx.WithContext(ctx).Create(reservation).Error
}

func (r *reservationRepository) FindByID(ctx context.Context, id string) (*entity.InventoryReservation, error) {
	var res entity.InventoryReservation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *entity.InventoryReservation) error {
	return r.db.WithContext(ctx).Save(reservation).Error
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status entity.ReservationStatus, updatedBy string) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).
		Model(&entity.InventoryReservation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_by": updatedBy,
		}).Error
}

// GetOverlappingReservations finds all active reservations that overlap with the given date range.
func (r *reservationRepository) GetOverlappingReservations(ctx context.Context, tx *gorm.DB, equipmentID string, startDate, endDate string) ([]entity.InventoryReservation, error) {
	var reservations []entity.InventoryReservation
	db := tx
	if db == nil {
		db = r.db
	}

	// Overlap condition: existing.start_date <= requested.end_date AND existing.end_date >= requested.start_date
	err := db.WithContext(ctx).
		Where("equipment_id = ?", equipmentID).
		Where("status IN ?", []string{
			string(entity.ReservationStatusPendingPayment),
			string(entity.ReservationStatusConfirmed),
			string(entity.ReservationStatusActive),
		}).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Find(&reservations).Error

	return reservations, err
}

// GetPeakUsage calculates the maximum concurrent reservation quantity for any single day
// within the specified date range using the reservation_date_locks table.
func (r *reservationRepository) GetPeakUsage(ctx context.Context, tx *gorm.DB, equipmentID string, startDate, endDate string) (uint, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var result struct {
		PeakUsage uint
	}

	err := db.WithContext(ctx).
		Raw(`SELECT COALESCE(MAX(daily_total), 0) AS peak_usage FROM (
			SELECT lock_date, SUM(reserved_qty) + SUM(maintenance_qty) AS daily_total
			FROM reservation_date_locks
			WHERE equipment_id = ?
			  AND lock_date BETWEEN ? AND ?
			GROUP BY lock_date
		) AS daily_totals`, equipmentID, startDate, endDate).
		Scan(&result).Error

	return result.PeakUsage, err
}

func (r *reservationRepository) CreateDateLocks(ctx context.Context, tx *gorm.DB, locks []entity.ReservationDateLock) error {
	if len(locks) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&locks).Error
}

func (r *reservationRepository) DeleteDateLocksByReservation(ctx context.Context, tx *gorm.DB, reservationID string) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).
		Where("reservation_id = ?", reservationID).
		Delete(&entity.ReservationDateLock{}).Error
}

func (r *reservationRepository) FindExpired(ctx context.Context, before time.Time, limit int) ([]entity.InventoryReservation, error) {
	var reservations []entity.InventoryReservation
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < ?", entity.ReservationStatusPendingPayment, before).
		Limit(limit).
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) ExpireReservation(ctx context.Context, tx *gorm.DB, id string) error {
	now := time.Now()
	return tx.WithContext(ctx).
		Model(&entity.InventoryReservation{}).
		Where("id = ? AND status = ?", id, entity.ReservationStatusPendingPayment).
		Updates(map[string]interface{}{
			"status":       entity.ReservationStatusExpired,
			"cancelled_at": now,
		}).Error
}

var _ repository.ReservationRepository = (*reservationRepository)(nil)
