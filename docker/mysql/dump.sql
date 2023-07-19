CREATE DATABASE IF NOT EXISTS `getground`;

USE `getground`;

--
-- Table structure for table `table`
--

DROP TABLE IF EXISTS `table`;

CREATE TABLE `table` (
  `id` int NOT NULL AUTO_INCREMENT,
  `capacity` int NOT NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;

--
-- Table structure for table `guest`
--

DROP TABLE IF EXISTS `guest`;

CREATE TABLE `guest` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET latin1 NOT NULL,
  `table` int NOT NULL,
  `accompanying_guests` int NOT NULL,
  `time_arrived` datetime NULL,
  PRIMARY KEY (`id`),
  KEY `guest_table_idx` (`table`),
  CONSTRAINT `guest_table` FOREIGN KEY (`table`) REFERENCES `table` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) DEFAULT CHARSET=utf8;
