package article

import (
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"

	"lmm/api/db"
)

type Category struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
}

func GetCategories(c *elesion.Context) {
	userIDStr := c.Params.ByName("user")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	categories, err := getCategories(userID)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("categories not found")
		return
	}
	c.Status(http.StatusOK).JSON(categories)
}

func getCategories(userID int64) ([]Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query(
		`SELECT id, user, name FROM category WHERE user = ? ORDER BY name`,
		userID,
	)
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

func GetArticleCategory(c *elesion.Context) {
	userIDStr := c.Params.ByName("user")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	articleIDStr := c.Params.ByName("article")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid article id: " + articleIDStr)
		return
	}

	category, err := getArticleCategory(userID, articleID)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("no such category")
		return
	}
	c.Status(http.StatusOK).JSON(category)
}

func getArticleCategory(userID, articleID int64) (*Category, error) {
	d := db.UseDefault()
	defer d.Close()

	var id int64
	err := d.QueryRow("SELECT category FROM article_category WHERE user = ? AND article = ?", userID, articleID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return getCategoryByID(id)
}

func getCategoryByID(id int64) (*Category, error) {
	d := db.UseDefault()
	defer d.Close()

	category := Category{}
	err := d.QueryRow("SELECT id, user, name FROM category WHERE id = ?", id).Scan(
		&category.ID, &category.User, &category.Name,
	)
	if err != nil {
		return nil, err
	}
	return &category, err
}

/*
func NewCategory(c *elesion.Context) {
	body := Category{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	_, err = newCategory(body)
	if err == nil {
		c.Status(http.StatusOK).String("success")
		return
	}
	if err == db.ErrAlreadyExists {
		c.Status(http.StatusConflict).String(body.Name + " already exists")
		return
	}
	c.Status(http.StatusInternalServerError).Error(err.Error()).String("unknown error")
}

func newCategory(category Category) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	// check if category already exists
	var id int64
	err := d.QueryRow("SELECT id FROM categories WHERE name = ?", category.Name).Scan(&id)
	if err == nil { // name already exists
		return id, db.ErrAlreadyExists
	}
	if err != sql.ErrNoRows {
		return id, err
	}
	// continue if no such row

	result, err := d.Exec("INSERT INTO categories (user_id, name) VALUES (?, ?)", category.UserID, category.Name)
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
