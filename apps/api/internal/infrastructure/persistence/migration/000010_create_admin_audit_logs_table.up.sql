-- Migration: 000010_create_admin_audit_logs_table
-- Description: Create admin_audit_logs table
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: admin_audit_logs
-- -----------------------------------------------------------
CREATE TABLE `admin_audit_logs` (
    `id`          CHAR(36) NOT NULL,
    `admin_id`    CHAR(36) NOT NULL,
    `action`      VARCHAR(50) NOT NULL,
    `entity_type` VARCHAR(50) NOT NULL,
    `entity_id`   CHAR(36) NOT NULL,
    `old_values`  JSON DEFAULT NULL,
    `new_values`  JSON DEFAULT NULL,
    `ip_address`  VARCHAR(45) DEFAULT NULL,
    `user_agent`  VARCHAR(500) DEFAULT NULL,
    `created_at`  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`  CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_audit_admin_created` (`admin_id`, `created_at`),
    INDEX `idx_audit_entity` (`entity_type`, `entity_id`),
    INDEX `idx_audit_action_created` (`action`, `created_at`),
    CONSTRAINT `fk_audit_admin` FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
