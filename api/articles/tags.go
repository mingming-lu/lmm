package articles

import (
	"lmm/api/db"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

type Tag struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func GetTags(c *elesion.Context) {
	userID := c.Query().Get("user_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	tags, err := getTags(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func getTags(userID string) ([]Tag, error) {
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
		itr.Scan(&tag.ID, &tag.UserID, &tag.Name)

		tags = append(tags, tag)
	}
	return tags, nil
}
