-- Migration: 000001_create_users_table
-- Description: Create users, user_sessions, and password_resets tables
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: users
-- -----------------------------------------------------------
CREATE TABLE `users` (
    `id`                CHAR(36) NOT NULL,
    `email`             VARCHAR(255) NOT NULL,
    `password_hash`     VARCHAR(255) NOT NULL,
    `full_name`         VARCHAR(100) NOT NULL,
    `phone`             VARCHAR(20) DEFAULT NULL,
    `avatar_url`        VARCHAR(500) DEFAULT NULL,
    `role`              ENUM('renter','owner','admin') NOT NULL DEFAULT 'renter',
    `email_verified_at` TIMESTAMP NULL DEFAULT NULL,
    `is_active`         TINYINT(1) NOT NULL DEFAULT 1,
    `created_at`        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`        TIMESTAMP NULL DEFAULT NULL,
    `created_by`        CHAR(36) DEFAULT NULL,
    `updated_by`        CHAR(36) DEFAULT NULL,
    `deleted_by`        CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_users_email` (`email`),
    INDEX `idx_users_role` (`role`),
    INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: user_sessions
-- -----------------------------------------------------------
CREATE TABLE `user_sessions` (
    `id`                 CHAR(36) NOT NULL,
    `user_id`            CHAR(36) NOT NULL,
    `refresh_token_hash` VARCHAR(255) NOT NULL,
    `user_agent`         VARCHAR(500) DEFAULT NULL,
    `ip_address`         VARCHAR(45) DEFAULT NULL,
    `expires_at`         TIMESTAMP NOT NULL,
    `revoked_at`         TIMESTAMP NULL DEFAULT NULL,
    `created_at`         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`         CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_sessions_user_id` (`user_id`),
    INDEX `idx_sessions_expires_at` (`expires_at`),
    CONSTRAINT `fk_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: password_resets
-- -----------------------------------------------------------
CREATE TABLE `password_resets` (
    `id`         CHAR(36) NOT NULL,
    `user_id`    CHAR(36) NOT NULL,
    `token_hash` VARCHAR(255) NOT NULL,
    `expires_at` TIMESTAMP NOT NULL,
    `used_at`    TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_password_resets_user_id` (`user_id`),
    CONSTRAINT `fk_password_resets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
