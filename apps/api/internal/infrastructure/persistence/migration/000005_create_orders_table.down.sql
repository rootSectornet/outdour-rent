-- Migration: 000005_create_orders_table (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `inventory_reservations` DROP FOREIGN KEY `fk_reservations_order_item`;
DROP TABLE IF EXISTS `order_status_history`;
DROP TABLE IF EXISTS `order_items`;
DROP TABLE IF EXISTS `orders`;

SET FOREIGN_KEY_CHECKS = 1;
