package articles

import (
	"lmm/api/db"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

type Category struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func GetCategories(c *elesion.Context) {
	userID := c.Query().Get("user_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	categories, err := getCategories(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(categories)
}

func getCategories(userID string) ([]Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user_id, name FROM categories WHERE user_id = ? ORDER BY name", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	categories := make([]Category, 0)
	for itr.Next() {
		category := Category{}
		itr.Scan(&category.ID, &category.UserID, &category.Name)

		categories = append(categories, category)
	}
	return categories, nil
}
