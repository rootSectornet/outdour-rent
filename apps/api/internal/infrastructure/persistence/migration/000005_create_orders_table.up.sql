-- Migration: 000005_create_orders_table
-- Description: Create orders, order_items, order_status_history tables
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: orders
-- -----------------------------------------------------------
CREATE TABLE `orders` (
    `id`                CHAR(36) NOT NULL,
    `order_number`      VARCHAR(20) NOT NULL,
    `renter_id`         CHAR(36) NOT NULL,
    `store_id`          CHAR(36) NOT NULL,
    `status`            ENUM('pending_payment','paid','approved','rejected','active','completed','cancelled','expired') NOT NULL,
    `rental_start_date` DATE NOT NULL,
    `rental_end_date`   DATE NOT NULL,
    `rental_days`       INT UNSIGNED NOT NULL,
    `subtotal`          DECIMAL(12,2) NOT NULL,
    `service_fee`       DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `deposit_total`     DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `total_amount`      DECIMAL(12,2) NOT NULL,
    `notes`             TEXT DEFAULT NULL,
    `pickup_method`     ENUM('pickup','delivery') NOT NULL DEFAULT 'pickup',
    `delivery_address`  TEXT DEFAULT NULL,
    `payment_deadline`  TIMESTAMP NOT NULL,
    `approved_at`       TIMESTAMP NULL DEFAULT NULL,
    `rejected_at`       TIMESTAMP NULL DEFAULT NULL,
    `rejection_reason`  VARCHAR(500) DEFAULT NULL,
    `completed_at`      TIMESTAMP NULL DEFAULT NULL,
    `cancelled_at`      TIMESTAMP NULL DEFAULT NULL,
    `created_at`        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`        TIMESTAMP NULL DEFAULT NULL,
    `created_by`        CHAR(36) DEFAULT NULL,
    `updated_by`        CHAR(36) DEFAULT NULL,
    `deleted_by`        CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_orders_order_number` (`order_number`),
    INDEX `idx_orders_renter_status` (`renter_id`, `status`),
    INDEX `idx_orders_store_status` (`store_id`, `status`),
    INDEX `idx_orders_status_deadline` (`status`, `payment_deadline`),
    INDEX `idx_orders_dates` (`rental_start_date`, `rental_end_date`),
    INDEX `idx_orders_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_orders_renter` FOREIGN KEY (`renter_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_orders_store` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_orders_dates` CHECK (`rental_end_date` >= `rental_start_date`),
    CONSTRAINT `chk_orders_amount` CHECK (`total_amount` >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: order_items
-- -----------------------------------------------------------
CREATE TABLE `order_items` (
    `id`               CHAR(36) NOT NULL,
    `order_id`         CHAR(36) NOT NULL,
    `equipment_id`     CHAR(36) NOT NULL,
    `reservation_id`   CHAR(36) NOT NULL,
    `quantity`         INT UNSIGNED NOT NULL,
    `price_per_day`    DECIMAL(12,2) NOT NULL,
    `rental_days`      INT UNSIGNED NOT NULL,
    `subtotal`         DECIMAL(12,2) NOT NULL,
    `deposit_per_unit` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `deposit_subtotal` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`       CHAR(36) DEFAULT NULL,
    `updated_by`       CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_order_items_order_id` (`order_id`),
    INDEX `idx_order_items_equipment_id` (`equipment_id`),
    UNIQUE INDEX `idx_order_items_reservation_id` (`reservation_id`),
    CONSTRAINT `fk_order_items_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_order_items_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_order_items_reservation` FOREIGN KEY (`reservation_id`) REFERENCES `inventory_reservations` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_order_items_qty` CHECK (`quantity` > 0),
    CONSTRAINT `chk_order_items_price` CHECK (`price_per_day` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add back-reference from inventory_reservations to order_items
ALTER TABLE `inventory_reservations`
    ADD CONSTRAINT `fk_reservations_order_item` FOREIGN KEY (`order_item_id`) REFERENCES `order_items` (`id`) ON DELETE SET NULL;

-- -----------------------------------------------------------
-- Table: order_status_history
-- -----------------------------------------------------------
CREATE TABLE `order_status_history` (
    `id`          CHAR(36) NOT NULL,
    `order_id`    CHAR(36) NOT NULL,
    `from_status` VARCHAR(30) DEFAULT NULL,
    `to_status`   VARCHAR(30) NOT NULL,
    `changed_by`  CHAR(36) DEFAULT NULL,
    `reason`      VARCHAR(500) DEFAULT NULL,
    `metadata`    JSON DEFAULT NULL,
    `created_at`  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`  CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_status_history_order_created` (`order_id`, `created_at`),
    CONSTRAINT `fk_status_history_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
