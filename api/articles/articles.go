package articles

import (
	"database/sql"
	"encoding/json"
	"lmm/api/db"
	"net/http"
	"strings"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"
)

type Response struct {
	Articles   []Article `json:"articles"`
	NextCursor string    `json:"next_cursor"`
}

type Result struct {
}

type Article struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
	CategoryID  int    `json:"category_id"`
}

func GetArticles(c *elesion.Context) {
	userID := c.Query().Get("user_id")
	categoryID := c.Query().Get("category_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	articles, err := getAllArticles(userID, categoryID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(articles)
}

func GetArticle(c *elesion.Context) {
	id := c.Query().Get("id")
	if id == "" {
		c.Status(http.StatusBadRequest).String("missing id")
		return
	}

	article, err := getArticle(id)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(article)
}

func getAllArticles(userID, categoryID string) ([]Article, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	var itr *sql.Rows
	var err error
	query := `SELECT id, title, text, created_date, updated_date, category_id FROM articles WHERE user_id = ? ORDER BY created_date DESC`
	if categoryID == "" {
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
		itr.Scan(&article.ID, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate, &article.CategoryID)

		articles = append(articles, article)
	}
	return articles, nil
}

func getArticle(id string) (*Article, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	article := Article{}
	err := d.QueryRow("SELECT id, title, text, created_date, updated_date, category_id FROM articles WHERE id = ?", id).Scan(
		&article.ID, &article.Title, &article.Text, &article.CreatedDate, &article.UpdatedDate, &article.CategoryID,
	)
	if err != nil {
		return nil, err
	}
	return &article, err
}

type postBody struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}

func PostArticle(c *elesion.Context) {
	// get user_id
	userID := c.Query().Get("user_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	// parse body (should be json)
	decoder := json.NewDecoder(c.Request.Body)

	body := postBody{}
	err := decoder.Decode(&body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}

	// insert into table
	_, err = postArticle(userID, body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func postArticle(userID string, body postBody) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	categoryID, err := getCategoryIDByName(body.Category)
	if err != nil {
		return -1, err
	}

	text := strings.TrimSpace(body.Text)

	result, err := d.Exec("INSERT INTO articles (user_id, title, text, category_id) VALUES (?, ?, ?, ?)", userID, body.Title, text, categoryID)
	if err != nil {
		return -1, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return -1, err
	} else if rows != 1 {
		return -1, errors.WithCaller("rows affected should be 1", 2)
	}
	return result.LastInsertId()
}
