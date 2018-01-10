package user

import "lmm/api/db"

const createUser = `
CREATE TABLE IF NOT EXISTS user (
	id int unsigned NOT NULL AUTO_INCREMENT,
	uid varchar(36) NOT NULL UNIQUE,
	token varchar(36) NOT NULL UNIQUE,
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

type User struct {
	ID          int64  `json:"id"`
	UID         int64  `json:"uid"`
	Token       string `json:"token"`
	CreatedDate string `json:"created_date"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
	Profession  string `json:"profession"`
	Location    string `json:"location"`
	Email       string `json:"email"`
}

func getUser(id int64) (*User, error) {
	d := db.Default()

}
