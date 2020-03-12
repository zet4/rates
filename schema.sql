CREATE TABLE `rates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `day` DATE NOT NULL,
  `currency` CHAR(3) NOT NULL,
  `rate` DECIMAL(10,5) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `day_currency` (`day`,`currency`)
);