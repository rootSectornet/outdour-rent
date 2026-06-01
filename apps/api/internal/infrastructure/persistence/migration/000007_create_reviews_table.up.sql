-- Migration: 000007_create_reviews_table
-- Description: Create reviews table
-- Created: 2026-06-01

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------
-- Table: reviews
-- -----------------------------------------------------------
CREATE TABLE `reviews` (
    `id`               CHAR(36) NOT NULL,
    `user_id`          CHAR(36) NOT NULL,
    `equipment_id`     CHAR(36) NOT NULL,
    `order_id`         CHAR(36) NOT NULL,
    `store_id`         CHAR(36) NOT NULL,
    `rating`           TINYINT UNSIGNED NOT NULL,
    `comment`          TEXT DEFAULT NULL,
    `owner_reply`      TEXT DEFAULT NULL,
    `owner_replied_at` TIMESTAMP NULL DEFAULT NULL,
    `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`       TIMESTAMP NULL DEFAULT NULL,
    `created_by`       CHAR(36) DEFAULT NULL,
    `updated_by`       CHAR(36) DEFAULT NULL,
    `deleted_by`       CHAR(36) DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_reviews_equipment_id` (`equipment_id`),
    INDEX `idx_reviews_store_id` (`store_id`),
    INDEX `idx_reviews_user_id` (`user_id`),
    UNIQUE INDEX `uidx_reviews_order_equipment` (`order_id`, `equipment_id`),
    INDEX `idx_reviews_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_reviews_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_reviews_equipment` FOREIGN KEY (`equipment_id`) REFERENCES `equipment` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_reviews_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_reviews_store` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `chk_reviews_rating` CHECK (`rating` >= 1 AND `rating` <= 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
