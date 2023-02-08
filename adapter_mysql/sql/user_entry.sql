CREATE TABLE `user_entry`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    bigint(20) unsigned NOT NULL,
    `feed_id`    bigint(20) unsigned NOT NULL,
    `entry_id`   bigint(20) unsigned NOT NULL,
    `readed`     tinyint(1) DEFAULT '0',
    `starred`    tinyint(1) DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_user_entry` (`user_id`,`entry_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户的 entry 数据