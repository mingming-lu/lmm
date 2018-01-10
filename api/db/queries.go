package db

const createUser = `
CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	uid varchar(32) NOT NULL UNIQUE,
	token varchar(32) NOT NULL UNIQUE,
	created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name varchar(32) NOT NULL,
	avatar_url varchar(255),
	description text,
	profession varchar(32),
	location varchar(255),
	email varchar(255),
	PRIMARY KEY (id)
)
`

const createSkill = `
CREATE TABLE IF NOT EXISTS skill (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id)
);
`

const createLanguage = `
CREATE TABLE IF NOT EXISTS language (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
	name varchar(32) NOT NULL UNIQUE,
	sort int unsigned NOT NULL UNIQUE,
	PRIMARY KEY (id)
);
`

const createEducation = `
CREATE TABLE IF NOT EXISTS education (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
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
	user_id int unsigned NOT NULL,
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
	user_id int unsigned NOT NULL,
	date date NOT NULL,
	name varchar(255) NOT NULL,
	PRIMARY KEY (id)
)
`

const createArticle = `
CREATE TABLE IF NOT EXISTS article (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
	created_date datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_date datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	category_id int unsigned NOT NULL DEFAULT 0,
	PRIMARY KEY (id)
)
`

const createArticleCategory = `
CREATE TABLE IF NOT EXISTS category (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (name)
)
`

const createArticleTag = `
CREATE TABLE IF NOT EXISTS tag (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_id int unsigned NOT NULL,
	article_id int unsigned NOT NULL,
	name varchar(32) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (name, article_id)
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
	createArticleCategory,
	createArticleTag,
}
