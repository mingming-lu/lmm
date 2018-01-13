package article

import (
	"encoding/json"
	"fmt"
	"lmm/api/db"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"
)

type Response struct {
	Articles   []Article `json:"article"`
	NextCursor string    `json:"next_cursor"`
}

type Article struct {
	ID          int64  `json:"id"`
	User        string `json:"user"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// GetArticles gets all articles according to user name or more information given by query parameters
// GET /users/:user/articles
func GetArticles(c *elesion.Context) {
	name := c.Params.ByName("user")

	values := c.Query()
	values.Set("user", name)
	articles, err := getArticles(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("article not found")
		return
	}
	c.Status(http.StatusOK).JSON(articles)
}

func getArticles(values url.Values) ([]Article, error) {
	d := db.UseDefault()
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT id, user, title, text, created_date, updated_date FROM article %s ORDER BY created_date DESC`,
		db.NewValues(values).Where(),
	)

	articles := make([]Article, 0)
	cursor, err := d.Query(query)
	if err != nil {
		return articles, err
	}
	defer cursor.Close()

	for cursor.Next() {
		article := Article{}
		err = cursor.Scan(&article.ID, &article.User, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate)
		if err != nil {
			return articles, err // return all articles found with error
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// GetArticle gets the article depending on user name and article id
// GET /user/:user/articles/:article
func GetArticle(c *elesion.Context) {
	user := c.Params.ByName("user")
	articleStr := c.Params.ByName("article")
	article, err := strconv.ParseInt(articleStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid id: " + articleStr)
		return
	}

	a, err := getArticle(user, article)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("article not found")
		return
	}
	c.Status(http.StatusOK).JSON(a)
}

func getArticle(user string, id int64) (*Article, error) {
	d := db.UseDefault()
	defer d.Close()

	article := Article{}
	err := d.QueryRow(
		"SELECT id, user, title, text, created_date, updated_date FROM article WHERE id = ? AND user = ?",
		id, user,
	).Scan(
		&article.ID, &article.User, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate,
	)
	if err != nil {
		return nil, err
	}
	return &article, err
}

// NewArticle post new article to the user given by url path
// POST /users/:user/articles
func NewArticle(c *elesion.Context) {
	user := c.Params.ByName("user")

	body := Article{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	body.User = user

	_, err = newArticle(body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String("failed to post article")
		return
	}

	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("success to post article but failed when post tags")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func newArticle(body Article) (int64, error) {
	d := db.UseDefault()
	defer d.Close()

	result, err := d.Exec(
		"INSERT INTO article (user, title, text) VALUES (?, ?, ?)",
		body.User, strings.TrimSpace(body.Title), strings.TrimSpace(body.Text),
	)
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

// UpdateArticle update the article where user name and article id are matched
// PUT /users/:user/articles/:article
func UpdateArticle(c *elesion.Context) {
	user := c.Params.ByName("user")
	articleIDStr := c.Params.ByName("article")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil || articleID <= 0 {
		c.Status(http.StatusBadRequest).String("invalid article id: " + articleIDStr)
		return
	}

	body := Article{}
	body.User = user
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	err = updateArticle(articleID, body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func updateArticle(id int64, article Article) error {
	d := db.UseDefault()
	defer d.Close()

	res, err := d.Exec("UPDATE article SET title = ?, text = ? WHERE id = ? AND user = ?",
		article.Title, article.Text, id, article.User,
	)
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.WithCaller("rows affected should be 1", 2)
	}
	return err
}

// DeleteArticle delete the article where user name and article id are matched
// DELETE /users/:user/articles/:article
func DeleteArticle(c *elesion.Context) {
	user := c.Params.ByName("user")
	articleStr := c.Params.ByName("article")
	article, err := strconv.ParseInt(articleStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid article id: " + articleStr)
		return
	}

	err = deleteArticle(user, article)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func deleteArticle(user string, id int64) error {
	d := db.UseDefault()
	defer d.Close()

	result, err := d.Exec("DELETE FROM article WHERE id = ? AND user = ?", id, user)
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.Newf("rows affected shoulb be 1 but got %d", rows)
	}
	return err
}
