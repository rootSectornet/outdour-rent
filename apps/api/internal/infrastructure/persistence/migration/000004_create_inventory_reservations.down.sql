-- Migration: 000004_create_inventory_reservations (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `reservation_date_locks`;
DROP TABLE IF EXISTS `inventory_reservations`;

SET FOREIGN_KEY_CHECKS = 1;
