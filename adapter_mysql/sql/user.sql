CREATE TABLE `user`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `username`   varchar(50)  NOT NULL,
    `hash`       varchar(100) NOT NULL,
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;