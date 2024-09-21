CREATE TABLE `prize_lists` (
  `prize_name` varchar(255) NOT NULL,
  `amount` decimal(20,8) NOT NULL,
  `probability` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`prize_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
