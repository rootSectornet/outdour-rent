-- Migration: 000004_create_inventory_reservations
-- Description: Create inventory_reservations and reservation_date_locks (core reservation engine)
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: inventory_reservations
-- Core of the date-range based rental reservation engine.
-- Status determines whether reserved qty counts toward availability.
-- Active statuses: pending_payment, confirmed, active
-- Terminal statuses: completed, cancelled, expired
-- -----------------------------------------------------------
CREATE TABLE `inventory_reservations` (
    `id`                  CHAR(36) NOT NULL,
    `equipment_id`        CHAR(36) NOT NULL,
    `order_item_id`       CHAR(36) DEFAULT NULL,
    `quantity`            INT UNSIGNED NOT NULL,
    `start_date`          DATE NOT NULL,
    `end_date`            DATE NOT NULL,
    `status`              ENUM('pending_payment','confirmed','active','completed','cancelled','expired') NOT NULL,
    `expires_at`          TIMESTAMP NULL DEFAULT NULL,
    `confirmed_at`        TIMESTAMP NULL DEFAULT NULL,
    `activated_at`        TIMESTAMP NULL DEFAULT NULL,
    `completed_at`        TIMESTAMP NULL DEFAULT NULL,
    `cancelled_at`        TIMESTAMP NULL DEFAULT NULL,
    `cancellation_reason` VARCHAR(500) DEFAULT NULL,
    `created_at`          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`          CHAR(36) DEFAULT NULL,
    `updated_by`          CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),

    -- Primary availability lookup: find active reservations for equipment in date range
    INDEX `idx_reservations_equip_status_dates` (`equipment_id`, `status`, `start_date`, `end_date`),

    -- Date range overlap queries
    INDEX `idx_reservations_equipment_dates` (`equipment_id`, `start_date`, `end_date`),

    -- Expiration cleanup job
    INDEX `idx_reservations_status_expires` (`status`, `expires_at`),

    -- Order item linkage
    UNIQUE INDEX `idx_reservations_order_item` (`order_item_id`),

    CONSTRAINT `fk_reservations_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_reservations_dates` CHECK (`end_date` >= `start_date`),
    CONSTRAINT `chk_reservations_qty` CHECK (`quantity` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: reservation_date_locks
-- Denormalized per-date quantity locks for fast O(1) availability.
-- Maintained transactionally with inventory_reservations.
-- One row per (equipment, date, reservation/maintenance).
-- Availability = total_stock - MAX(SUM per date)
-- -----------------------------------------------------------
CREATE TABLE `reservation_date_locks` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `equipment_id`    CHAR(36) NOT NULL,
    `lock_date`       DATE NOT NULL,
    `reserved_qty`    INT UNSIGNED NOT NULL DEFAULT 0,
    `maintenance_qty` INT UNSIGNED NOT NULL DEFAULT 0,
    `reservation_id`  CHAR(36) DEFAULT NULL,
    `maintenance_id`  CHAR(36) DEFAULT NULL,
    `created_at`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),

    -- Prevent duplicate lock entries
    UNIQUE INDEX `uidx_date_locks_equip_date_res_maint` (`equipment_id`, `lock_date`, `reservation_id`, `maintenance_id`),

    -- Fast availability aggregation: SUM(reserved_qty + maintenance_qty) GROUP BY lock_date
    INDEX `idx_date_locks_availability` (`equipment_id`, `lock_date`),

    CONSTRAINT `fk_date_locks_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_date_locks_reservation` FOREIGN KEY (`reservation_id`) REFERENCES `inventory_reservations` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_date_locks_maintenance` FOREIGN KEY (`maintenance_id`) REFERENCES `equipment_maintenance` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
