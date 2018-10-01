CREATE TABLE IF NOT EXISTS `user` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	`password` VARCHAR(255) NOT NULL,
	`token` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `name` (`name`),
	UNIQUE `token` (`token`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `article` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` BIGINT UNSIGNED NOT NULL, -- user.id
	`uid` VARCHAR(255) NOT NULL,
	`title` VARCHAR(255) NOT NULL,
	`body` TEXT NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	INDEX `created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `article_tag` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`article` INT UNSIGNED NOT NULL, -- article.id
	`sort` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `tag_id` (`article`, `sort`),
	INDEX `tag_name` (`article`, `sort`, `name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `image` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`uid` VARCHAR(255) NOT NULL,
	`user` BIGINT UNSIGNED NOT NULL,
	`type` TINYINT UNSIGNED NOT NULL DEFAULT 0,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	INDEX `created_at` (`type`, `created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
