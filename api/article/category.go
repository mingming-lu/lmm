package article

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"

	"lmm/api/db"
	"lmm/api/user"
)

type Category struct {
	ID      int64  `json:"id"`
	User    int64  `json:"user"`
	Article int64  `json:"article"`
	Name    string `json:"name"`
}

// GET /categories
func GetCategories(c *elesion.Context) {
	queryParams := c.Query()
	if len(queryParams) == 0 {
		c.Status(http.StatusBadRequest).String("empty query parameter")
		return
	}

	values := db.NewValuesFromURL(queryParams)
	categories, err := getCategories(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("categories not found")
		return
	}
	c.Status(http.StatusOK).JSON(categories)
}

func getCategories(values db.Values) ([]Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT MIN(id), user, name FROM category %s GROUP BY name ORDER BY name`,
		values.Where(),
	)

	itr, err := d.Query(query)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	categories := make([]Category, 0)
	for itr.Next() {
		category := Category{}
		err = itr.Scan(&category.ID, &category.User, &category.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}
	return categories, nil
}

func NewCategory(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	body := Category{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}

	if body.Article == 0 || body.Name == "" {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("empty name or article")
		return
	}

	body.User = usr.ID
	id, err := newCategory(body)
	switch err {
	case nil:
		c.Status(http.StatusOK).Header("Location", fmt.Sprintf("/categories?id=%d", id)).String("success")
	case db.ErrAlreadyExists:
		c.Status(http.StatusConflict).String("category already exists: " + body.Name)
	default:
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
	}
}

func newCategory(category Category) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	ok, err := d.Exists("SELECT 1 FROM category WHERE name = ?", category.Name)
	if err != nil {
		return 0, err
	}
	if ok {
		return 0, db.ErrAlreadyExists
	}

	result, err := d.Exec("INSERT INTO category (user, article, name) VALUES (?, ?, ?)",
		category.User, category.Article, category.Name,
	)
	if err != nil {
		return 0, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return 0, err
	} else if rows != 1 {
		return 0, errors.Newf("rows affected are expected to be 1 but got %d", 2)
	}

	return result.LastInsertId()
}

/*
func UpdateCategory(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
	}

	body := Category{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	err = updateCategory(id, body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid input")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func updateCategory(id int64, body Category) error {
	d := db.New().Use("lmm")
	defer d.Close()

	result, err := d.Exec(
		"UPDATE categories SET name = ? WHERE id = ? AND user_id = ?",
		body.Name, id, body.UserID,
	)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.WithCaller("rows affected should be 1", 2)
	}
	return nil
}

func DeleteCategory(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	err = deleteCategory(id)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("not exists id: " + idStr)
		return
	}
	c.Status(http.StatusOK).String("success")
}

func deleteCategory(id int64) error {
	d := db.New().Use("lmm")
	defer d.Close()

	result, err := d.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.Newf("rows affected should be 1 but got", rows)
	}
	return nil
}
*/
