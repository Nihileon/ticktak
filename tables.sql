create table if not exists t_user
(
    `id`          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `username`    VARCHAR(64),
    `password`    VARCHAR(64),
    `description` text,
    `create_time` timestamp        NOT NULL DEFAULT '1999-01-01 00:00:00',
    `modify_time` timestamp        NOT NULL DEFAULT '1999-01-01 00:00:00',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'user table';


create table if not exists t_task
(
    `id`          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `username`    VARCHAR(64)      NOT NULL,
    `title`       VARCHAR(256)     NOT NULL,
    `state`       TINYINT          NOT NULL DEFAULT 0,
    `priority`    TINYINT          Not NULL DEFAULT 0,
    `content`     text,
    `create_time` timestamp        NOT NULL DEFAULT '1999-01-01 00:00:00',
    `modify_time` timestamp        NOT NULL DEFAULT '1999-01-01 00:00:00',
    PRIMARY KEY (`id`),
    KEY `idx_username` (`username`),
    KEY `idx_username_title` (`username`, `title`),
    KEY `idx_username_priority` (`username`,`priority`),
    KEY `idx_username_state` (`username`,`state`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'item table';

