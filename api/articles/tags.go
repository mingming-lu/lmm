package articles

import (
	"lmm/api/db"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

type Tag struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
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

	itr, err := d.Query("SELECT id, user_id, name FROM tags WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	tags := make([]Tag, 0)
	for itr.Next() {
		tag := Tag{}
		err = itr.Scan(&tag.ID, &tag.UserID, &tag.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}
