-- Migration: 000009_create_mountains_table (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `mountain_equipment_recs`;
DROP TABLE IF EXISTS `mountains`;

SET FOREIGN_KEY_CHECKS = 1;
