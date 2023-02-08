drop table if exists `user_feed`;

CREATE TABLE `user_feed`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    bigint(20) unsigned NOT NULL,
    `feed_id`    bigint(20) unsigned NOT NULL,
    `tag_name`   varchar(256) NOT NULL default 'default',
    `title`      varchar(256) NOT NULL default '',
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_user_feed` (`user_id`, `feed_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户订阅的 feed
--
