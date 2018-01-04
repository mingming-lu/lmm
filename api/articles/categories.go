package articles

import (
	"database/sql"
	"encoding/json"
	"lmm/api/db"
	"net/http"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"
)

type Category struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func GetCategories(c *elesion.Context) {
	userID := c.Query().Get("user_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	categories, err := getCategories(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(categories)
}

func getCategories(userID string) ([]Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query(
		`SELECT id, user_id, name FROM categories WHERE user_id = ? ORDER BY name`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	categories := make([]Category, 0)
	for itr.Next() {
		category := Category{}
		itr.Scan(&category.ID, &category.UserID, &category.Name)

		categories = append(categories, category)
	}
	return categories, nil
}

func NewCategory(c *elesion.Context) {
	body := Category{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}

	_, err = newCategory(body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func newCategory(body Category) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	// check if category already exists
	var id int64
	err := d.QueryRow("SELECT id FROM categories WHERE name = ?", body.Name).Scan(&id)
	if err == nil {
		return id, errors.New("category " + body.Name + "already exists")
	}
	if err != sql.ErrNoRows {
		return id, err
	}
	// continue if no such row

	result, err := d.Exec("INSERT INTO categories (user_id, name) VALUES (?, ?)", body.UserID, body.Name)
	if err != nil {
		return 0, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return 0, err
	} else if rows != 1 {
		return 0, errors.WithCaller("rows affected should be 1", 2)
	}

	return result.LastInsertId()
}
