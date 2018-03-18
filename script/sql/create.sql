CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	name varchar(32) NOT NULL UNIQUE,
	password varchar(128) NOT NULL,
	guid varchar(36) NOT NULL UNIQUE,
	token varchar(36) NOT NULL UNIQUE,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS blog (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE (user, title)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS blog_category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	blog int unsigned NOT NULL,
	category int unsigned NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(blog),
	INDEX `category` (category)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS tag (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	blog int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, blog, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS image (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(255) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS photo (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	image int unsigned NOT NULL,
	deleted tinyint unsigned NOT NULL DEFAULT 0,
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
