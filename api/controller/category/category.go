package category

import (
	"encoding/json"
	"fmt"
	"lmm/api/db"
	model "lmm/api/domain/model/category"
	usecase "lmm/api/usecase/category"
	"lmm/api/usecase/user"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

func Register(c *elesion.Context) {
	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	m := model.Minimal{}
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	// TODO handle already exists
	id, err := usecase.Register(usr.ID, blogID, m.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}

	c.Header("Location", fmt.Sprintf("categories/%d", id)).Status(http.StatusCreated).String("success")
}

func Update(c *elesion.Context) {
	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	category, err := usecase.FetchByBlog(blogID)
	if err != nil {
		c.Status(http.StatusNotFound).String("Blog not found").Error(err.Error())
		return
	}

	if category.User != usr.ID {
		c.Status(http.StatusForbidden).String("Access forbidden").Error(err.Error())
		return
	}

	m := model.Minimal{}
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	err = usecase.Update(usr.ID, blogID, m.Name)
	if err != nil && err == db.ErrNoChange {
		c.Status(http.StatusNoContent).String("No change").Error(err.Error())
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}

	c.Status(http.StatusOK).String("success")
}

func GetByBlog(c *elesion.Context) {
	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
	}

	blog, err := usecase.FetchByBlog(blogID)
	if err != nil {
		c.Status(http.StatusNotFound).String("Not Found").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(blog)
}

func GetByUser(c *elesion.Context) {
	blogID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
	}

	blog, err := usecase.FetchByUser(blogID)
	if err != nil {
		c.Status(http.StatusNotFound).String("Not Found").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(blog)
}
