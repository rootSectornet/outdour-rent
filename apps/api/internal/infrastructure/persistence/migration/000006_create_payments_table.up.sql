-- Migration: 000006_create_payments_table
-- Description: Create payments, deposits, refunds tables
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: payments
-- -----------------------------------------------------------
CREATE TABLE `payments` (
    `id`                       CHAR(36) NOT NULL,
    `order_id`                 CHAR(36) NOT NULL,
    `payment_type`             ENUM('rental','deposit') NOT NULL,
    `amount`                   DECIMAL(12,2) NOT NULL,
    `method`                   VARCHAR(50) DEFAULT NULL,
    `status`                   ENUM('pending','success','failed','expired','refunded') NOT NULL,
    `midtrans_transaction_id`  VARCHAR(100) DEFAULT NULL,
    `midtrans_order_id`        VARCHAR(100) DEFAULT NULL,
    `midtrans_snap_token`      VARCHAR(255) DEFAULT NULL,
    `midtrans_redirect_url`    VARCHAR(500) DEFAULT NULL,
    `midtrans_payment_type`    VARCHAR(50) DEFAULT NULL,
    `paid_at`                  TIMESTAMP NULL DEFAULT NULL,
    `expired_at`               TIMESTAMP NULL DEFAULT NULL,
    `raw_callback`             JSON DEFAULT NULL,
    `created_at`               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`               CHAR(36) DEFAULT NULL,
    `updated_by`               CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_payments_order_type` (`order_id`, `payment_type`),
    INDEX `idx_payments_midtrans_order` (`midtrans_order_id`),
    INDEX `idx_payments_status` (`status`),
    CONSTRAINT `fk_payments_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_payments_amount` CHECK (`amount` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: deposits
-- -----------------------------------------------------------
CREATE TABLE `deposits` (
    `id`               CHAR(36) NOT NULL,
    `order_id`         CHAR(36) NOT NULL,
    `payment_id`       CHAR(36) DEFAULT NULL,
    `amount`           DECIMAL(12,2) NOT NULL,
    `status`           ENUM('pending','held','returned','forfeited','partially_returned') NOT NULL,
    `returned_amount`  DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `forfeited_amount` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `forfeit_reason`   VARCHAR(500) DEFAULT NULL,
    `returned_at`      TIMESTAMP NULL DEFAULT NULL,
    `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`       CHAR(36) DEFAULT NULL,
    `updated_by`       CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_deposits_order_id` (`order_id`),
    INDEX `idx_deposits_status` (`status`),
    CONSTRAINT `fk_deposits_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_deposits_payment` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE SET NULL,
    CONSTRAINT `chk_deposits_amount` CHECK (`amount` > 0),
    CONSTRAINT `chk_deposits_returned` CHECK (`returned_amount` <= `amount`),
    CONSTRAINT `chk_deposits_forfeited` CHECK (`forfeited_amount` <= `amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
-- Table: refunds
-- -----------------------------------------------------------
CREATE TABLE `refunds` (
    `id`                 CHAR(36) NOT NULL,
    `payment_id`         CHAR(36) NOT NULL,
    `order_id`           CHAR(36) NOT NULL,
    `amount`             DECIMAL(12,2) NOT NULL,
    `reason`             VARCHAR(500) NOT NULL,
    `status`             ENUM('pending','processing','success','failed') NOT NULL,
    `midtrans_refund_id` VARCHAR(100) DEFAULT NULL,
    `processed_at`       TIMESTAMP NULL DEFAULT NULL,
    `created_at`         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`         CHAR(36) DEFAULT NULL,
    `updated_by`         CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_refunds_payment_id` (`payment_id`),
    INDEX `idx_refunds_order_id` (`order_id`),
    CONSTRAINT `fk_refunds_payment` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_refunds_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_refunds_amount` CHECK (`amount` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
