DROP DATABASE IF EXISTS `lmm_test`;
CREATE DATABASE IF NOT EXISTS `lmm_test`;

USE `lmm_test`;

CREATE TABLE IF NOT EXISTS `user` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	`password` VARCHAR(255) NOT NULL,
	`token` VARCHAR(255) NOT NULL,
	`role` VARCHAR(31) NOT NULL,
	`created_at` DATETIME NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `name` (`name`),
	UNIQUE `token` (`token`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `user_role_change_history` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`operator` BIGINT UNSIGNED NOT NULL, -- user.id
	`operator_role` VARCHAR(31) NOT NULL, -- user.role
	`target_user` BIGINT UNSIGNED NOT NULL, -- user.id
	`from_role` VARCHAR(31) NOT NULL, -- user.role
	`to_role` VARCHAR(31) NOT NULL, -- user.role
	`changed_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	INDEX `changed_at` (`changed_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `user_password_change_history` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` BIGINT UNSIGNED NOT NULL, -- user.id
	`changed_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	INDEX `user_change_history` (`user`, `changed_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `article` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` BIGINT UNSIGNED NOT NULL, -- user.id
	`uid` VARCHAR(255) NOT NULL,
	`alias_uid` VARCHAR(255) NOT NULL,
	`title` VARCHAR(255) NOT NULL,
	`body` TEXT NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	UNIQUE `alias_uid` (`alias_uid`),
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

CREATE TABLE IF NOT EXISTS `asset` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	`type` TINYINT UNSIGNED NOT NULL, -- image:0 photo:1
	`user` INT UNSIGNED NOT NULL, -- user.id
	`created_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `name` (`name`),
	INDEX `created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `image_alt` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`asset` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `name` (`asset`, `name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
