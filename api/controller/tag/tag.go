package tag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lmm/api/usecase/user"

	model "lmm/api/domain/model/tag"
	usecase "lmm/api/usecase/tag"

	"github.com/akinaru-lu/elesion"
)

func Register(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	tag := model.Minimal{}
	err = json.NewDecoder(c.Request.Body).Decode(&tag)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	tagID, err := usecase.Add(usr.ID, blogID, tag.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Header("Location", fmt.Sprintf("/v1/blog/%d/tags/%d", blogID, tagID)).Status(http.StatusCreated).String("success")
}

func Update(c *elesion.Context) {
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

	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	tagID, err := strconv.ParseInt(c.Params.ByName("tag"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid tag id").Error(err.Error())
		return
	}

	err = usecase.Update(usr.ID, blogID, tagID, m.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func GetByUser(c *elesion.Context) {
	userID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user id").Error(err.Error())
		return
	}

	tags, err := usecase.FetchByUser(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func GetByBlog(c *elesion.Context) {
	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	tags, err := usecase.FetchByBlog(blogID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func Delete(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	blogID, err := strconv.ParseInt(c.Params.ByName("blog"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid blog id").Error(err.Error())
		return
	}

	tagID, err := strconv.ParseInt(c.Params.ByName("tag"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid tag id").Error(err.Error())
		return
	}

	err = usecase.Delete(usr.ID, blogID, tagID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
