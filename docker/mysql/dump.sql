CREATE DATABASE IF NOT EXISTS `getground`;

USE `getground`;

--
-- Table structure for table `table`
--

DROP TABLE IF EXISTS `table`;

CREATE TABLE `table` (
  `id` int NOT NULL AUTO_INCREMENT,
  `capacity` int NOT NULL,
  `reserved_seats` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;

--
-- Table structure for table `guest`
--

DROP TABLE IF EXISTS `guest`;

CREATE TABLE `guest` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL UNIQUE,
  `table_id` int NOT NULL,
  `accompanying_guests` int NOT NULL,
  `time_arrived` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  KEY `guest_table_idx` (`table_id`),
  CONSTRAINT `guest_table` FOREIGN KEY (`table_id`) REFERENCES `table` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) DEFAULT CHARSET=utf8;
