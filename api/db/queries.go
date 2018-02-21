package db

const createUser = `
CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	name varchar(32) NOT NULL UNIQUE,
	password varchar(128) NOT NULL,
	guid varchar(36) NOT NULL UNIQUE,
	token varchar(36) NOT NULL UNIQUE,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createProfile = `
CREATE TABLE IF NOT EXISTS profile (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL UNIQUE,
	nickname varchar(32) NOT NULL UNIQUE,
	avatar_url varchar(128) NOT NULL DEFAULT '',
	description varchar(400) NOT NULL DEFAULT '',
	profession varchar(32) NOT NULL DEFAULT '',
	location varchar(128) NOT NULL DEFAULT '',
	email varchar(128) NOT NULL DEFAULT '',
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createSkill = `
CREATE TABLE IF NOT EXISTS skill (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createLanguage = `
CREATE TABLE IF NOT EXISTS language (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createEducation = `
CREATE TABLE IF NOT EXISTS education (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	date_from date NOT NULL,
	date_to date NOT NULL,
	institution varchar(255) NOT NULL,
	department varchar(255),
	major varchar(255),
	degree varchar(32) NOT NULL,
	current bit NOT NULL DEFAULT b'0',
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createWorkExperience = `
CREATE TABLE IF NOT EXISTS work_experience (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	date_from date NOT NULL,
	date_to date NOT NULL,
	company varchar(255) NOT NULL,
	position varchar(32),
	status varchar(32),
	current bit NOT NULL DEFAULT b'0',
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createQualification = `
CREATE TABLE IF NOT EXISTS qualification (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	date date NOT NULL,
	name varchar(255) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createBlog = `
CREATE TABLE IF NOT EXISTS blog (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createCategory = `
CREATE TABLE IF NOT EXISTS category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	blog int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, blog)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createTag = `
CREATE TABLE IF NOT EXISTS tag (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	blog int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, blog, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

const createImage = `
CREATE TABLE IF NOT EXISTS image (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	type varchar (15) NOT NULL,
	url varchar(127) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE (url)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`

var CreateSQL = []string{
	createUser,
	createProfile,
	createSkill,
	createLanguage,
	createEducation,
	createWorkExperience,
	createQualification,
	createBlog,
	createCategory,
	createTag,
	createImage,
}
