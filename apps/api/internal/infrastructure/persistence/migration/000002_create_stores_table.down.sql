-- Migration: 000002_create_stores_table (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `store_documents`;
DROP TABLE IF EXISTS `stores`;

SET FOREIGN_KEY_CHECKS = 1;
