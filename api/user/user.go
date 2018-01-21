package user

import (
	"encoding/json"
	"fmt"
	"lmm/api/db"
	"lmm/api/utils"
	"lmm/api/utils/token"
	"net/http"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"
)

type User struct {
	ID          int64
	Name        string
	Password    string
	GUID        string
	Token       string
	CreatedDate string
}

type UserSmall struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func GetUser(values db.Values) (*User, error) {
	d := db.UseDefault()
	defer d.Close()

	user := User{}
	query := fmt.Sprintf("SELECT id, name, guid, token, created_at FROM user %s", values.Where())
	err := d.QueryRow(query).Scan(&user.ID, &user.Name, &user.GUID, &user.Token, &user.CreatedDate)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
	user.Password = utils.ToBase64([]byte(user.GUID + user.Password))

	result, err := d.Exec(`INSERT INTO user (name, password, guid, token) VALUES (?, ?, ?, ?)`,
		user.Name, user.Password, user.GUID, user.Token,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func SignUp(c *elesion.Context) {
	user := User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid input").Error(err.Error())
		return
	}

	if user.Name == "" || user.Password == "" {
		c.Status(http.StatusBadRequest).String("empty input")
		return
	}

	id, err := newUser(user)
	if err != nil {
		c.Status(http.StatusInternalServerError).String(err.Error()).Error(err.Error())
		return
	}
	c.Writer.Header().Set("Location", fmt.Sprintf("/users/%d/login", id))
	c.Status(http.StatusCreated).String("success")
}

func Login(c *elesion.Context) {
	user := User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body").Error(err.Error())
		return
	}
	res, err := fetchUser(user.Name, user.Password)
	if err != nil {
		c.Status(http.StatusNotFound).String(err.Error()).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(res)
}

func fetchUser(name, password string) (*UserSmall, error) {
	d := db.UseDefault()
	defer d.Close()

	values := db.NewValues()
	values["name"] = name
	user, err := GetUser(values)
	if err != nil {
		return nil, err
	}

	pwEncoded := utils.ToBase64([]byte(user.GUID + password))
	ok, err := d.Exists("SELECT 1 FROM user WHERE password = ? LIMIT 1", pwEncoded)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("password error")
	}

	accessToken, err := token.Encode(user.Token)
	if err != nil {
		return nil, errors.Wrap(err, "refresh token failed")
	}

	return &UserSmall{ID: user.ID, Name: user.Name, Token: accessToken}, nil
}

func Logout(c *elesion.Context) {
	Verify(c)
}

func Verify(c *elesion.Context) {
	accessToken := c.Request.Header.Get("Authorization")
	originToken, err := token.Decode(accessToken)
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token")
		return
	}

	values := db.NewValues()
	values["token"] = originToken
	_, err = GetUser(values)
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token")
		return
	}
	c.Status(http.StatusOK).String("success")
}

// NewTestUser create a user for testing
/*
func NewTestUser() (*User, error) {
	user := User{}
	user.Name = "test"
	user.Password = "test password"
	id, err := newUser(user)
	if err != nil {
		panic(err)
	}

	values := db.NewValues()
	values["id"] = id
	values["name"] = user.Name
	return GetUser(values)
}
*/
