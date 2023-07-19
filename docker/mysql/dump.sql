CREATE DATABASE  IF NOT EXISTS `getground`;

USE `getground`;

--
-- Table structure for table `guest`
--

DROP TABLE IF EXISTS `guest`;

CREATE TABLE `guest` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) CHARACTER SET latin1 NOT NULL,
  `table` int(11) NOT NULL,
  `accompanying_guests` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `guest_table_idx` (`table`),
  CONSTRAINT `guest_table` FOREIGN KEY (`table`) REFERENCES `table` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Table structure for table `table`
--

DROP TABLE IF EXISTS `table`;

CREATE TABLE `table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `capacity` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
