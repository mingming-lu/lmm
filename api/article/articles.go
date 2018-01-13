package article

import (
	"database/sql"
	"encoding/json"
	"lmm/api/db"
	"net/http"
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
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
	CategoryID  int64  `json:"category_id"`
	Tags        []Tag  `json:"tags"`
}

func GetArticles(c *elesion.Context) {
	userIDStr := c.Params.ByName("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	categoryIDStr := c.Query().Get("category_id")
	if categoryIDStr == "" {
		categoryIDStr = "0"
	}
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil || categoryID < 0 { // allow zero ID (for all categories)
		c.Status(http.StatusBadRequest).String("invalid category id: " + categoryIDStr)
		return
	}

	articles, err := getArticles(userID, categoryID)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("article not found")
		return
	}
	c.Status(http.StatusOK).JSON(articles)
}

func getArticles(userID, categoryID int64) ([]Article, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	var itr *sql.Rows
	var err error
	query := `SELECT id, title, text, created_date, updated_date, category_id FROM article WHERE user_id = ? ORDER BY created_date DESC`
	if categoryID == 0 {
		itr, err = d.Query(query, userID)
	} else {
		itr, err = d.Query(query+` AND category_id = ?`, userID, categoryID)
	}
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	articles := make([]Article, 0)
	for itr.Next() {
		article := Article{}
		err = itr.Scan(&article.ID, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate, &article.CategoryID)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	return articles, nil
}

func GetArticle(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	article, err := getArticle(id)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("article not found")
		return
	}
	c.Status(http.StatusOK).JSON(article)
}

func getArticle(id int64) (*Article, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	article := Article{}
	err := d.QueryRow(
		"SELECT id, title, text, created_date, updated_date, category_id FROM article WHERE id = ?", id).Scan(
		&article.ID, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate, &article.CategoryID,
	)
	if err != nil {
		return nil, err
	}
	return &article, err
}

func NewArticle(c *elesion.Context) {
	// parse body (should be json)
	body := Article{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	// insert into table
	id, err := postArticle(body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("failed to post article")
		return
	}

	for i := range body.Tags {
		body.Tags[i].UserID = body.UserID
		body.Tags[i].ArticleID = id
	}
	_, err = newTags(body.Tags)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("success to post article but failed when post tags")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func postArticle(body Article) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	result, err := d.Exec(
		"INSERT INTO article (user_id, title, text, category_id) VALUES (?, ?, ?, ?)",
		body.UserID, body.Title, strings.TrimSpace(body.Text), body.CategoryID,
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

func UpdateArticle(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	body := Article{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}

	err = updateArticle(id, body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid input")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func updateArticle(id int64, body Article) error {
	d := db.New().Use("lmm")
	defer d.Close()

	res, err := d.Exec("UPDATE article SET title = ?, text = ?, category_id = ? WHERE id = ? AND user_id = ?",
		body.Title, body.Text, body.CategoryID, id, body.UserID,
	)
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.WithCaller("rows affected should be 1", 2)
	}
	return err
}
