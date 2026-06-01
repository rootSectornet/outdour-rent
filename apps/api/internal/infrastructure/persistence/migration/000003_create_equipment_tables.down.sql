-- Migration: 000003_create_equipment_tables (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `equipment_maintenance`;
DROP TABLE IF EXISTS `equipment_pricing`;
DROP TABLE IF EXISTS `equipment_photos`;
DROP TABLE IF EXISTS `equipment`;
DROP TABLE IF EXISTS `equipment_categories`;

SET FOREIGN_KEY_CHECKS = 1;
