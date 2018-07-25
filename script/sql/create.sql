CREATE TABLE IF NOT EXISTS `user` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(31) NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`token` VARCHAR(63) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE `name` (`name`),
	UNIQUE `token` (`token`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `blog` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` BIGINT UNSIGNED NOT NULL,
	`title` VARCHAR(63) NOT NULL,
	`text` TEXT NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `title` (`title`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `category` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(31) NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `name` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `blog_category` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`blog` INT UNSIGNED NOT NULL,
	`category` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `blog` (`blog`),
	INDEX `category` (`category`, `blog`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `tag` (
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
	`blog` INT unsigned NOT NULL,
	`name` VARCHAR(31) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `blog_tag` (`blog`, `name`),
	INDEX `tag` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `image` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`uid` VARCHAR(255) NOT NULL,
	`user` INT UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	INDEX `created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS photo (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	image int unsigned NOT NULL,
	deleted tinyint unsigned NOT NULL DEFAULT 0,
	last_modified timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE `image` (`image`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS project (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(63) NOT NULL,
	icon varchar(255) NOT NULL DEFAULT "",
	url varchar(255) NOT NULL DEFAULT "",
	description varchar(1023) NOT NULL DEFAULT "",
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	from_date date,
	to_date date,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
