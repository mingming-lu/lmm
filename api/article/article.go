package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"

	"lmm/api/db"
	"lmm/api/user"
)

type Article struct {
	ID        int64  `json:"id"`
	User      int64  `json:"user"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetArticles gets all articles according to query parameters
// GET /articles
func GetArticles(c *elesion.Context) {
	queryParams := c.Query()
	if len(queryParams) == 0 {
		c.Status(http.StatusBadRequest).String("empty query parameter")
		return
	}

	values := db.NewValuesFromURL(c.Query())
	articles, err := getArticles(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("article not found")
		return
	}
	c.Status(http.StatusOK).JSON(articles)
}

func getArticles(values db.Values) ([]Article, error) {
	d := db.UseDefault()
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT id, user, title, text, created_at, updated_at FROM article %s ORDER BY created_at DESC`,
		values.Where(),
	)

	articles := make([]Article, 0)
	cursor, err := d.Query(query)
	if err != nil {
		return articles, err
	}
	defer cursor.Close()

	for cursor.Next() {
		article := Article{}
		err = cursor.Scan(&article.ID, &article.User, &article.Title, &article.Text, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return articles, err // return all articles found with error
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// NewArticle post new article to the user given by url path
// POST /articles
func NewArticle(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	body := Article{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	body.User = usr.ID

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

// NewTestArticle creates a new user, and creates a new article by the created user
/*
func NewTestArticle() (*Article, *user.UserProfile) {
	usr := user.NewTestUser()
	id, err := newArticle(Article{
		ID:    usr.ID,
		Title: "test",
		Text:  "This is a test article",
	})
	if err != nil {
		panic(err)
	}
	article, err := getArticle(usr.ID, id)
	if err != nil {
		panic(err)
	}
	return article, usr
}
*/

// UpdateArticle update the article where user name and article id are matched
// PUT /articles
func UpdateArticle(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	body := Article{}
	body.User = usr.ID
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	err = updateArticle(body)

	switch err {
	case nil:
		c.Status(http.StatusOK).String("success")
	case db.ErrNoChange:
		c.Status(http.StatusAccepted).String("no change")
	case db.ErrNoRows:
		c.Status(http.StatusNotFound).String(fmt.Sprintf("no such article: %d", body.ID))
	default:
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
	}
}

func updateArticle(article Article) error {
	d := db.UseDefault()
	defer d.Close()

	ok, err := d.Exists("SELECT 1 FROM article WHERE id = ? AND user = ?", article.ID, article.User)
	if err != nil {
		return err
	}
	if !ok {
		return db.ErrNoRows
	}

	res, err := d.Exec("UPDATE article SET title = ?, text = ? WHERE id = ? AND user = ?",
		article.Title, article.Text, article.ID, article.User,
	)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	} else if rows > 1 {
		return errors.Newf("rows affected should be larger than 1 but got ", rows)
	}
	return nil
}

// DeleteArticle delete the article where user name and article id are matched
// DELETE /articles
func DeleteArticle(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	article := Article{}
	err = json.NewDecoder(c.Request.Body).Decode(&article)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invald body")
		return
	}

	err = deleteArticle(usr.ID, article.ID)
	switch err {
	case nil:
		c.Status(http.StatusOK).String("success")
	case db.ErrNoRows:
		c.Status(http.StatusNotFound).Stringf("article not found")
	default:
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
	}
}

func deleteArticle(user, id int64) error {
	d := db.UseDefault()
	defer d.Close()

	ok, err := d.Exists("SELECT 1 FROM article WHERE id = ? AND user = ?", id, user)
	if err != nil {
		return err
	}
	if !ok {
		return db.ErrNoRows
	}

	result, err := d.Exec("DELETE FROM article WHERE id = ? AND user = ?", id, user)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	} else if rows > 1 {
		return errors.Newf("rows affected shoulb be 1 but got %d", rows)
	}
	return nil
}
