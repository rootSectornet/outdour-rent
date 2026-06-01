-- Rollback: 000011_add_google_oauth_to_users

DROP INDEX `idx_users_google_id` ON `users`;

ALTER TABLE `users`
    DROP COLUMN `provider`,
    DROP COLUMN `google_id`,
    MODIFY COLUMN `password_hash` VARCHAR(255) NOT NULL;
