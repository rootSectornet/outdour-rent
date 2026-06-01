-- Revert status and purchase_date columns from equipment table
ALTER TABLE `equipment`
    DROP INDEX `idx_equipment_status`,
    DROP COLUMN `purchase_date`,
    DROP COLUMN `status`;
