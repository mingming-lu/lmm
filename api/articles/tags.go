package articles

import (
	"encoding/json"
	"lmm/api/db"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"
)

type Tag struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ArticleID int64  `json:"article_id"`
	Name      string `json:"name"`
}

func GetTags(c *elesion.Context) {
	userIDStr := c.Params.ByName("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	tags, err := getTags(userID)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("tags not found")
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func getTags(userID int64) ([]Tag, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT name FROM tags WHERE user_id = ? GROUP BY name ORDER BY name", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	tags := make([]Tag, 0)
	for itr.Next() {
		tag := Tag{}
		err = itr.Scan(&tag.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

func NewTags(c *elesion.Context) {
	var tags []Tag
	err := json.NewDecoder(c.Request.Body).Decode(&tags)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	newTags(tags)
	c.Status(http.StatusInternalServerError).Error("not implemented")
}

func newTags(tags []Tag) (int64, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	query := "INSERT INTO tags (user_id, article_id, name) VALUES "
	var values []interface{}
	for _, tag := range tags {
		query += "(?, ?, ?), "
		values = append(values, tag.UserID, tag.ArticleID, tag.Name)
	}
	query += "ON DUPLICATE KEY UPDATE"

	stmtIns, err := d.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(values)
	if err != nil {
		return 0, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return 0, err
	} else if rows != int64(len(tags)) {
		return 0, errors.Newf("rows affected should be %d, but got %d", len(tags), rows)
	}
	return result.LastInsertId()
}
