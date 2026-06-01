-- Migration: 000006_create_payments_table (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `refunds`;
DROP TABLE IF EXISTS `deposits`;
DROP TABLE IF EXISTS `payments`;

SET FOREIGN_KEY_CHECKS = 1;
