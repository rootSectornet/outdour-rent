-- Migration: 000011_add_google_oauth_to_users
-- Description: Add Google OAuth support columns to users table
-- Created: 2026-06-01

ALTER TABLE `users`
    ADD COLUMN `google_id` VARCHAR(255) DEFAULT NULL AFTER `is_active`,
    ADD COLUMN `provider` VARCHAR(50) NOT NULL DEFAULT 'local' AFTER `google_id`,
    MODIFY COLUMN `password_hash` VARCHAR(255) DEFAULT NULL;

CREATE UNIQUE INDEX `idx_users_google_id` ON `users` (`google_id`);
