package user

import (
	"encoding/json"
	"fmt"
	"github.com/akinaru-lu/elesion"
	"lmm/api/db"
	"lmm/api/utils"
	"net/http"
	"strconv"
)

type User struct {
	ID          int64  `json:"id"`
	GUID        string `json:"guid"`
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

// GET /users/:user
// user: user id
func GetUser(c *elesion.Context) {
	idStr := c.Params.ByName("user")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	u, err := getUser(id)
	if err != nil {
		c.Status(http.StatusNotFound).String("user not found")
		return
	}
	c.Status(http.StatusOK).JSON(u)
}

func getUser(id int64) (*User, error) {
	d := db.UseDefault()
	defer d.Close()

	u := User{}
	err := d.QueryRow(
		"SELECT id, guid, token, created_date, name, nickname, avatar_url, description, profession, location, email FROM user WHERE id = ?", id,
	).Scan(
		&u.ID, &u.GUID, &u.Token, &u.CreatedDate, &u.Name, &u.Nickname, &u.AvatarURL, &u.Description, &u.Profession, &u.Location, &u.Email,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// POST /users
// body : name, nickname
func NewUser(c *elesion.Context) {
	usr := User{}
	err := json.NewDecoder(c.Request.Body).Decode(&usr)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	if usr.Name == "" || usr.Nickname == "" {
		c.Status(http.StatusBadRequest).String("empty name or nickname")
		return
	}

	id, err := newUser(usr)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String("invalid input")
		return
	}
	c.Writer.Header().Set("Location", fmt.Sprintf("/users/%d", id))
	c.Status(http.StatusCreated).JSON(newUser)
}

func newUser(user User) (int64, error) {
	d := db.UseDefault()
	defer d.Close()

	if user.GUID == "" {
		user.GUID = utils.NewUUID()
	}
	if user.Token == "" {
		user.Token = utils.NewUUID()
	}

	result, err := d.Exec(`INSERT INTO user (guid, token, name, nickname) VALUES (?, ?, ?, ?)`, user.GUID, user.Token, user.Name, user.Nickname)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
