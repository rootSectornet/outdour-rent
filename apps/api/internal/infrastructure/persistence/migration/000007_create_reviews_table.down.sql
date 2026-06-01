-- Migration: 000007_create_reviews_table (DOWN)

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `reviews`;

SET FOREIGN_KEY_CHECKS = 1;
