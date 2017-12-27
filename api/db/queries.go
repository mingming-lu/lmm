package db

const createProfile = `
CREATE TABLE IF NOT EXISTS profile (
	id int NOT NULL AUTO_INCREMENT,
	name varchar(255) NOT NULL,
	avatar_url varchar(255),
	bio text,
	location varchar(255),
	profession varchar(32),
	email varchar(255),
	PRIMARY KEY (id)
);
`

const createSkill = `
CREATE TABLE IF NOT EXISTS skill (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	name varchar(32) NOT NULL,
	level varchar(32) NOT NULL,
	sort tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (sort)
);
`

const createLanguage = `
CREATE TABLE IF NOT EXISTS language (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	name varchar(32) NOT NULL,
	level varchar(32) NOT NULL,
	sort tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (sort)
);
`

const createEducation = `
CREATE TABLE IF NOT EXISTS education (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	date_from date NOT NULL,
	date_to date NOT NULL,
	institution varchar(255) NOT NULL,
	department varchar(255),
	major varchar(255),
	degree varchar(32) NOT NULL,
	current bit NOT NULL,
	sort tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (sort)
);
`

const createWorkExperience = `
CREATE TABLE IF NOT EXISTS work_experience (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	date_from date NOT NULL,
	date_to date NOT NULL,
	company varchar(255) NOT NULL,
	position varchar(32),
	status varchar(32),
	current bit NOT NULL,
	sort tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (sort)
)
`

const createQualification = `
CREATE TABLE IF NOT EXISTS qualification (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	date date NOT NULL,
	name varchar(255) NOT NULL,
	sort tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (sort)
)
`

const createArticle = `
CREATE TABLE IF NOT EXISTS article (
	id int NOT NULL AUTO_INCREMENT,
	user_id int NOT NULL,
	created_date date NOT NULL,
	edited_date date NOT NULL,
	title varchar(255) NOT NULL,
	text text NOT NULL,
	PRIMARY KEY (id)
)
`

var CreateSQL = []string{
	createProfile,
	createSkill,
	createLanguage,
	createEducation,
	createWorkExperience,
	createQualification,
	createArticle,
}
