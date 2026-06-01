-- Migration: 000002_create_stores_table
-- Description: Create stores and store_documents tables
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: stores
-- -----------------------------------------------------------
CREATE TABLE `stores` (
    `id`           CHAR(36) NOT NULL,
    `owner_id`     CHAR(36) NOT NULL,
    `name`         VARCHAR(100) NOT NULL,
    `slug`         VARCHAR(120) NOT NULL,
    `description`  TEXT DEFAULT NULL,
    `phone`        VARCHAR(20) NOT NULL,
    `email`        VARCHAR(255) NOT NULL,
    `address`      TEXT NOT NULL,
    `city`         VARCHAR(100) NOT NULL,
    `province`     VARCHAR(100) NOT NULL,
    `postal_code`  VARCHAR(10) DEFAULT NULL,
    `latitude`     DECIMAL(10,8) DEFAULT NULL,
    `longitude`    DECIMAL(11,8) DEFAULT NULL,
    `logo_url`     VARCHAR(500) DEFAULT NULL,
    `banner_url`   VARCHAR(500) DEFAULT NULL,
    `status`       ENUM('pending','verified','suspended','rejected') NOT NULL DEFAULT 'pending',
    `verified_at`  TIMESTAMP NULL DEFAULT NULL,
    `rating_avg`   DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `rating_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `created_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`   TIMESTAMP NULL DEFAULT NULL,
    `created_by`   CHAR(36) DEFAULT NULL,
    `updated_by`   CHAR(36) DEFAULT NULL,
    `deleted_by`   CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_stores_slug` (`slug`),
    INDEX `idx_stores_owner_id` (`owner_id`),
    INDEX `idx_stores_status` (`status`),
    INDEX `idx_stores_city` (`city`),
    INDEX `idx_stores_location` (`latitude`, `longitude`),
    INDEX `idx_stores_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_stores_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: store_documents
-- -----------------------------------------------------------
CREATE TABLE `store_documents` (
    `id`              CHAR(36) NOT NULL,
    `store_id`        CHAR(36) NOT NULL,
    `document_type`   ENUM('ktp','npwp','siu','other') NOT NULL,
    `document_url`    VARCHAR(500) NOT NULL,
    `verified_at`     TIMESTAMP NULL DEFAULT NULL,
    `rejected_reason` VARCHAR(500) DEFAULT NULL,
    `created_at`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_store_docs_store_id` (`store_id`),
    CONSTRAINT `fk_store_docs_store` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
