package user

import (
	"encoding/json"
	"github.com/akinaru-lu/elesion"
	"lmm/api/db"
	"lmm/api/utils"
	"net/http"
	"strconv"
)

type User struct {
	ID          int64  `json:"id"`
	UID         string `json:"uid"`
	Token       string `json:"token"`
	CreatedDate string `json:"created_date"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
	Profession  string `json:"profession"`
	Location    string `json:"location"`
	Email       string `json:"email"`
}

func GetUser(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	user, err := getUser(id)
	if err != nil {
		c.Status(http.StatusNotFound).String("user not found")
		return
	}
	c.Status(http.StatusOK).JSON(user)
}

func getUser(id int64) (*User, error) {
	d := db.UseDefault()
	defer d.Close()

	user := User{}
	err := d.QueryRow(
		"SELECT id, uid, token, created_date, name, nickname, avatar_url, description, profession, location, email FROM user WHERE id = ?", id,
	).Scan(
		&user.ID, &user.UID, &user.Token, &user.CreatedDate, &user.Name, &user.Nickname, &user.AvatarURL, &user.Description, &user.Profession, &user.Location, &user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUser(c *elesion.Context) {
	user := User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body")
		return
	}

	if user.Name == "" || user.Nickname == "" {
		c.Status(http.StatusBadRequest).String("empty name or nickname")
		return
	}

	id, err := newUser(user)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid input")
		return
	}
	newUser, err := getUser(id)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String("something is wrong")
		return
	}
	c.Status(http.StatusCreated).JSON(newUser)
}

func newUser(user User) (int64, error) {
	d := db.UseDefault()
	defer d.Close()

	if user.UID == "" {
		user.UID = utils.NewUUID()
	}
	if user.Token == "" {
		user.Token = utils.NewUUID()
	}

	result, err := d.Exec(`INSERT INTO user (uid, token, name, nickname) VALUES (?, ?, ?, ?)`, user.UID, user.Token, user.Name, user.Nickname)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
