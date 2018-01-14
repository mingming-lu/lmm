package db

const createUser = `
CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	guid varchar(36) NOT NULL UNIQUE,
	token varchar(36) NOT NULL UNIQUE,
	created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name varchar(32) NOT NULL,
	nickname varchar(32) NOT NULL UNIQUE,
	avatar_url varchar(127) NOT NULL DEFAULT '',
	description varchar(400) NOT NULL DEFAULT '',
	profession varchar(32) NOT NULL DEFAULT '',
	location varchar(127) NOT NULL DEFAULT '',
	email varchar(127) NOT NULL DEFAULT '',
	PRIMARY KEY (id)
)
`

const createSkill = `
CREATE TABLE IF NOT EXISTS skill (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id),
	UNIQUE (user, name)
);
`

const createLanguage = `
CREATE TABLE IF NOT EXISTS language (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id),
	UNIQUE (user, name)
);
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
);
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
)
`

const createQualification = `
CREATE TABLE IF NOT EXISTS qualification (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	date date NOT NULL,
	name varchar(255) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, name)
)
`

const createArticle = `
CREATE TABLE IF NOT EXISTS article (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user int unsigned NOT NULL,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	created_date datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_date datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
)
`

const createCategory = `
CREATE TABLE IF NOT EXISTS category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user varchar(32) NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, name)
)
`

const createTag = `
CREATE TABLE IF NOT EXISTS tag (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user varchar(32) NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (user, name)
)
`

var CreateSQL = []string{
	createUser,
	createSkill,
	createLanguage,
	createEducation,
	createWorkExperience,
	createQualification,
	createArticle,
	createCategory,
	createTag,
}
