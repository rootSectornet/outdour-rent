-- Migration: 000003_create_equipment_tables
-- Description: Create equipment_categories, equipment, equipment_photos, equipment_pricing, equipment_maintenance
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: equipment_categories
-- -----------------------------------------------------------
CREATE TABLE `equipment_categories` (
    `id`         CHAR(36) NOT NULL,
    `parent_id`  CHAR(36) DEFAULT NULL,
    `name`       VARCHAR(100) NOT NULL,
    `slug`       VARCHAR(120) NOT NULL,
    `icon_url`   VARCHAR(500) DEFAULT NULL,
    `sort_order` INT NOT NULL DEFAULT 0,
    `is_active`  TINYINT(1) NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `created_by` CHAR(36) DEFAULT NULL,
    `updated_by` CHAR(36) DEFAULT NULL,
    `deleted_by` CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_categories_slug` (`slug`),
    INDEX `idx_categories_parent_id` (`parent_id`),
    INDEX `idx_categories_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_categories_parent` FOREIGN KEY (`parent_id`) REFERENCES `equipment_categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: equipment
-- -----------------------------------------------------------
CREATE TABLE `equipment` (
    `id`               CHAR(36) NOT NULL,
    `store_id`         CHAR(36) NOT NULL,
    `category_id`      CHAR(36) NOT NULL,
    `name`             VARCHAR(200) NOT NULL,
    `slug`             VARCHAR(220) NOT NULL,
    `description`      TEXT DEFAULT NULL,
    `brand`            VARCHAR(100) DEFAULT NULL,
    `specifications`   JSON DEFAULT NULL,
    `total_stock`      INT UNSIGNED NOT NULL DEFAULT 0,
    `condition`        ENUM('new','good','fair','worn') NOT NULL DEFAULT 'good',
    `weight_grams`     INT UNSIGNED DEFAULT NULL,
    `min_rental_days`  INT UNSIGNED NOT NULL DEFAULT 1,
    `max_rental_days`  INT UNSIGNED NOT NULL DEFAULT 30,
    `requires_deposit` TINYINT(1) NOT NULL DEFAULT 1,
    `deposit_amount`   DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `is_active`        TINYINT(1) NOT NULL DEFAULT 1,
    `rating_avg`       DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `rating_count`     INT UNSIGNED NOT NULL DEFAULT 0,
    `rental_count`     INT UNSIGNED NOT NULL DEFAULT 0,
    `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`       TIMESTAMP NULL DEFAULT NULL,
    `created_by`       CHAR(36) DEFAULT NULL,
    `updated_by`       CHAR(36) DEFAULT NULL,
    `deleted_by`       CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_equipment_store_slug` (`store_id`, `slug`),
    INDEX `idx_equipment_store_id` (`store_id`),
    INDEX `idx_equipment_category_id` (`category_id`),
    INDEX `idx_equipment_active` (`is_active`, `deleted_at`),
    INDEX `idx_equipment_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_equipment_store` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_equipment_category` FOREIGN KEY (`category_id`) REFERENCES `equipment_categories` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_equipment_stock` CHECK (`total_stock` >= 0),
    CONSTRAINT `chk_equipment_rental_days` CHECK (`max_rental_days` >= `min_rental_days`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: equipment_photos
-- -----------------------------------------------------------
CREATE TABLE `equipment_photos` (
    `id`            CHAR(36) NOT NULL,
    `equipment_id`  CHAR(36) NOT NULL,
    `photo_url`     VARCHAR(500) NOT NULL,
    `thumbnail_url` VARCHAR(500) DEFAULT NULL,
    `sort_order`    INT NOT NULL DEFAULT 0,
    `is_primary`    TINYINT(1) NOT NULL DEFAULT 0,
    `created_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`    CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_photos_equipment_sort` (`equipment_id`, `sort_order`),
    CONSTRAINT `fk_photos_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: equipment_pricing
-- -----------------------------------------------------------
CREATE TABLE `equipment_pricing` (
    `id`            CHAR(36) NOT NULL,
    `equipment_id`  CHAR(36) NOT NULL,
    `pricing_type`  ENUM('daily','weekly','monthly','custom') NOT NULL,
    `min_days`      INT UNSIGNED NOT NULL DEFAULT 1,
    `max_days`      INT UNSIGNED DEFAULT NULL,
    `price_per_day` DECIMAL(12,2) NOT NULL,
    `is_active`     TINYINT(1) NOT NULL DEFAULT 1,
    `created_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`    CHAR(36) DEFAULT NULL,
    `updated_by`    CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_pricing_equipment_active` (`equipment_id`, `is_active`),
    INDEX `idx_pricing_type` (`equipment_id`, `pricing_type`),
    CONSTRAINT `fk_pricing_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE CASCADE,
    CONSTRAINT `chk_pricing_price` CHECK (`price_per_day` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: equipment_maintenance
-- -----------------------------------------------------------
CREATE TABLE `equipment_maintenance` (
    `id`           CHAR(36) NOT NULL,
    `equipment_id` CHAR(36) NOT NULL,
    `quantity`     INT UNSIGNED NOT NULL,
    `start_date`   DATE NOT NULL,
    `end_date`     DATE NOT NULL,
    `reason`       VARCHAR(500) DEFAULT NULL,
    `created_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`   CHAR(36) DEFAULT NULL,
    `updated_by`   CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_maintenance_equipment_dates` (`equipment_id`, `start_date`, `end_date`),
    INDEX `idx_maintenance_daterange` (`start_date`, `end_date`, `equipment_id`),
    CONSTRAINT `fk_maintenance_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE CASCADE,
    CONSTRAINT `chk_maintenance_dates` CHECK (`end_date` >= `start_date`),
    CONSTRAINT `chk_maintenance_qty` CHECK (`quantity` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
