-- Rollback: 000012_create_store_photos_and_hours

DROP TABLE IF EXISTS `store_operating_hours`;
DROP TABLE IF EXISTS `store_photos`;

ALTER TABLE `stores`
    DROP COLUMN `suspended_at`,
    MODIFY COLUMN `status` ENUM('pending','verified','suspended','rejected') NOT NULL DEFAULT 'pending';

UPDATE `stores` SET `status` = 'pending' WHERE `status` = 'pending_approval';
UPDATE `stores` SET `status` = 'verified' WHERE `status` = 'active';
