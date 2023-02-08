CREATE TABLE `feed`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(256) NOT NULL default '',
    `feed_url`   varchar(256) NOT NULL default '',
    `home_url`   TEXT NULL,
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_feed_url` (`feed_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- feed