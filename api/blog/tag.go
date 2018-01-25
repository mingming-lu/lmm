package blog

import (
	"fmt"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/db"
)

type Tag struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
}

// GetTags is the handler of GET /users/:user/tags
// get all tags under the given user (id)
func GetTags(c *elesion.Context) {
	queryParams := c.Query()
	if len(queryParams) == 0 {
		c.Status(http.StatusBadRequest).String("empty query parameter")
		return
	}

	values := db.NewValuesFromURL(queryParams)

	tags, err := getTags(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("tags not found")
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func getTags(values db.Values) ([]Tag, error) {
	d := db.UseDefault()
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT MIN(id), user, name FROM tag %s GROUP BY name ORDER BY name`,
		values.Where(),
	)

	itr, err := d.Query(query)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	tags := make([]Tag, 0)
	for itr.Next() {
		tag := Tag{}
		err = itr.Scan(&tag.ID, &tag.User, &tag.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}
