-- Add status and purchase_date columns to equipment table
ALTER TABLE `equipment`
    ADD COLUMN `status` ENUM('available', 'reserved', 'rented', 'maintenance', 'damaged', 'retired')
        NOT NULL DEFAULT 'available' AFTER `condition`,
    ADD COLUMN `purchase_date` DATE NULL AFTER `status`,
    ADD INDEX `idx_equipment_status` (`status`);
