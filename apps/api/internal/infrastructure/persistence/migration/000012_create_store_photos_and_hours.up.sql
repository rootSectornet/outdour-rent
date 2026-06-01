-- Migration: 000012_create_store_photos_and_hours
-- Description: Create store_photos and store_operating_hours tables, update store status enum
-- Created: 2026-06-01

-- Update store status enum
ALTER TABLE `stores`
    MODIFY COLUMN `status` ENUM('pending_approval','active','suspended','rejected') NOT NULL DEFAULT 'pending_approval',
    ADD COLUMN `suspended_at` TIMESTAMP NULL DEFAULT NULL AFTER `verified_at`;

-- Migrate existing data
UPDATE `stores` SET `status` = 'pending_approval' WHERE `status` = 'pending';
UPDATE `stores` SET `status` = 'active' WHERE `status` = 'verified';

-- -----------------------------------------------------------
-- Table: store_photos
-- -----------------------------------------------------------
CREATE TABLE `store_photos` (
    `id`            CHAR(36) NOT NULL,
    `store_id`      CHAR(36) NOT NULL,
    `photo_url`     VARCHAR(500) NOT NULL,
    `thumbnail_url` VARCHAR(500) DEFAULT NULL,
    `caption`       VARCHAR(255) DEFAULT NULL,
    `sort_order`    INT NOT NULL DEFAULT 0,
    `is_primary`    TINYINT(1) NOT NULL DEFAULT 0,
    `created_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`    CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_store_photos_store_id` (`store_id`),
    CONSTRAINT `fk_store_photos_store` FOREIGN KEY (`store_id`) REFERENCES `stores`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: store_operating_hours
-- -----------------------------------------------------------
CREATE TABLE `store_operating_hours` (
    `id`          CHAR(36) NOT NULL,
    `store_id`    CHAR(36) NOT NULL,
    `day_of_week` TINYINT UNSIGNED NOT NULL COMMENT '1=Monday, 7=Sunday',
    `open_time`   VARCHAR(5) NOT NULL COMMENT 'HH:MM',
    `close_time`  VARCHAR(5) NOT NULL COMMENT 'HH:MM',
    `is_closed`   TINYINT(1) NOT NULL DEFAULT 0,
    `created_at`  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`  CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_store_hours_store_day` (`store_id`, `day_of_week`),
    CONSTRAINT `fk_store_hours_store` FOREIGN KEY (`store_id`) REFERENCES `stores`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
