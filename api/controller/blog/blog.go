package blog

import (
	"encoding/json"
	"fmt"
	model "lmm/api/domain/model/blog"
	usecase "lmm/api/usecase/blog"
	"lmm/api/usecase/user"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

func Post(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	blog := model.Minimal{}
	err = json.NewDecoder(c.Request.Body).Decode(&blog)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body").Error(err.Error())
		return
	}

	id, err := usecase.Post(usr.ID, blog.Title, blog.Text)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Header("Location", fmt.Sprintf("/blogs/%d", id)).Status(http.StatusCreated).String("success")
}

func Get(c *elesion.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	blog, err := usecase.Fetch(id)
	if err != nil {
		c.Status(http.StatusNotFound).String("No such blog").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(blog)
}

func GetByUser(c *elesion.Context) {
	userID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user ID").Error(err.Error())
		return
	}

	blogs, err := usecase.FetchByUser(userID)
	if err != nil {
		c.Status(http.StatusNotFound).String("Blogs not found").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(blogs)
}

func Update(c *elesion.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	m := model.Minimal{}
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body").Error(err.Error())
		return
	}

	_, err = usecase.Update(id, m.Title, m.Text)
	if err != nil {
		c.Status(http.StatusNotFound).String("No such blog").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func Delete(c *elesion.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	_, err = usecase.Delete(id)
	if err != nil {
		c.Status(http.StatusNotFound).String("No such blog").Error(err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
