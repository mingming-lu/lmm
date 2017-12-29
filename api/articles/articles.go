package articles

import (
	"lmm/api/db"
	"net/http"

	"github.com/akinaru-lu/elesion"
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
	EditedTime  string `json:"edited_date"`
}

func Handler(c *elesion.Context) {
	articles, err := allArticlesByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
	}
	c.Status(http.StatusOK).JSON(articles)
}

func allArticlesByUserID(id int) ([]Article, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, title, text, created_date, edited_date FROM articles WHERE user_id = ?", 1)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	articles := make([]Article, 0)
	for itr.Next() {
		article := Article{}
		itr.Scan(&article.ID, &article.Title, &article.Text, &article.CreatedDate, &article.EditedTime)

		articles = append(articles, article)
	}
	return articles, nil
}
