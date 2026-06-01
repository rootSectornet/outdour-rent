-- Migration: 000008_create_notifications_table
-- Description: Create notifications table
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: notifications
-- -----------------------------------------------------------
CREATE TABLE `notifications` (
    `id`         CHAR(36) NOT NULL,
    `user_id`    CHAR(36) NOT NULL,
    `type`       VARCHAR(50) NOT NULL,
    `title`      VARCHAR(200) NOT NULL,
    `body`       TEXT NOT NULL,
    `data`       JSON DEFAULT NULL,
    `channel`    ENUM('in_app','email','push') NOT NULL DEFAULT 'in_app',
    `read_at`    TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_notifications_user_read` (`user_id`, `read_at`),
    INDEX `idx_notifications_user_type` (`user_id`, `type`, `created_at`),
    CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
