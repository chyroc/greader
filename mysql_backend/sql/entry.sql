CREATE TABLE `entry`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `feed_id`    bigint(20) unsigned NOT NULL,
    `title`      varchar(256) NOT NULL,
    `url`        varchar(256) NOT NULL,
    `author`     varchar(256) NOT NULL default '',
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_feed_url` (`feed_id`, `url`),
    INDEX        `idx_feed_id` (`feed_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 存放 fetch 的文章数据