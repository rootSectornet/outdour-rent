-- Migration: 000009_create_mountains_table
-- Description: Create mountains and mountain_equipment_recs tables
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: mountains
-- -----------------------------------------------------------
CREATE TABLE `mountains` (
    `id`               CHAR(36) NOT NULL,
    `name`             VARCHAR(100) NOT NULL,
    `slug`             VARCHAR(120) NOT NULL,
    `description`      TEXT DEFAULT NULL,
    `elevation_meters` INT UNSIGNED DEFAULT NULL,
    `difficulty`       ENUM('easy','moderate','hard','expert') NOT NULL,
    `location`         VARCHAR(200) DEFAULT NULL,
    `latitude`         DECIMAL(10,8) DEFAULT NULL,
    `longitude`        DECIMAL(11,8) DEFAULT NULL,
    `image_url`        VARCHAR(500) DEFAULT NULL,
    `is_active`        TINYINT(1) NOT NULL DEFAULT 1,
    `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`       CHAR(36) DEFAULT NULL,
    `updated_by`       CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_mountains_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: mountain_equipment_recs
-- -----------------------------------------------------------
CREATE TABLE `mountain_equipment_recs` (
    `id`          CHAR(36) NOT NULL,
    `mountain_id` CHAR(36) NOT NULL,
    `category_id` CHAR(36) NOT NULL,
    `priority`    ENUM('essential','recommended','optional') NOT NULL,
    `notes`       VARCHAR(500) DEFAULT NULL,
    `created_at`  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`  CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_recs_mountain_priority` (`mountain_id`, `priority`),
    UNIQUE INDEX `uidx_recs_mountain_category` (`mountain_id`, `category_id`),
    CONSTRAINT `fk_recs_mountain` FOREIGN KEY (`mountain_id`) REFERENCES `mountains` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_recs_category` FOREIGN KEY (`category_id`) REFERENCES `equipment_categories` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
