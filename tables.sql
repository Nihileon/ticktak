create table if not exists t_user(
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64),
  `password` VARCHAR(64),
  `description` text,
  `create_time` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
  `modify_time` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_username` (`user`),
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'user table';
