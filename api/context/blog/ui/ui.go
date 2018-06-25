package ui

import (
	"fmt"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/service"
	"lmm/api/http"
	"lmm/api/storage"
	"log"
)

type UI struct {
	app *appservice.AppService
}

func New(db *storage.DB) *UI {
	app := appservice.New(db)
	return &UI{app: app}
}

func (ui *UI) PostBlog(c *http.Context) {
	user := c.Values().Get("user").(*account.User)
	blogID, err := ui.app.PostNewBlog(user, c.Request.Body)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/blog/%d", blogID)).String(http.StatusCreated, "success")
	case service.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, service.ErrEmptyBlogTitle.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) GetAllBlog(c *http.Context) {
	blogPage, err := ui.app.GetBlogListByPage(
		c.Request.Query("count"),
		c.Request.Query("page"),
	)

	switch err {
	case nil:
		c.JSON(http.StatusOK, blogPage)
	case service.ErrInvalidCount, service.ErrInvalidPage:
		c.String(http.StatusBadRequest, err.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) GetBlog(c *http.Context) {
	blog, err := ui.app.GetBlogByID(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, blog)
	case service.ErrNoSuchBlog:
		c.String(http.StatusNotFound, service.ErrNoSuchBlog.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) UpdateBlog(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	err := ui.app.EditBlog(user, c.Request.Path.Params("blog"), c.Request.Body)
	switch err {
	case nil:
		c.String(http.StatusOK, "success")
	case service.ErrBlogNoChange:
		http.NoContent(c)
	case service.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, service.ErrEmptyBlogTitle.Error())
	case service.ErrNoPermission:
		c.String(http.StatusForbidden, service.ErrNoSuchBlog.Error())
	case service.ErrNoSuchBlog:
		c.String(http.StatusNotFound, service.ErrNoSuchBlog.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) SetBlogCategory(c *http.Context) {
	_, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	err := ui.app.SetBlogCategory(c.Request.Path.Params("blog"), c.Request.Body)
	switch err {
	case nil:
		c.String(http.StatusOK, "success")
	case service.ErrNoSuchBlog, service.ErrNoSuchCategory:
		c.String(http.StatusBadRequest, err.Error())
	default:
		http.InternalServerError(c)
	}
}

func (ui *UI) PostCategory(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	categoryID, err := ui.app.RegisterNewCategory(user, c.Request.Body)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/categories/%d", categoryID)).String(http.StatusCreated, "success")
	case service.ErrInvalidCategoryName, service.ErrDuplicateCategoryName:
		c.String(http.StatusBadRequest, err.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) UpdateCategory(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	err := ui.app.EditCategory(user, c.Request.Path.Params("category"), c.Request.Body)
	switch err {
	case nil:
		c.String(http.StatusOK, "success")
	case service.ErrCategoryNoChanged:
		http.NoContent(c)
	case service.ErrInvalidCategoryName:
		c.String(http.StatusBadRequest, service.ErrInvalidCategoryName.Error())
	case service.ErrNoSuchCategory:
		c.String(http.StatusNotFound, service.ErrNoSuchCategory.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) GetAllCategoris(c *http.Context) {
	categories, err := ui.app.GetAllCategories()
	switch err {
	case nil:
		c.JSON(http.StatusOK, categories)
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

func (ui *UI) GetBlogCagetory(c *http.Context) {
	category, err := ui.app.GetCategoryOfBlog(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, category)
	case service.ErrNoSuchBlog:
		c.String(http.StatusNotFound, service.ErrNoSuchBlog.Error())
	case service.ErrCategoryNotSet:
		c.String(http.StatusNotFound, service.ErrCategoryNotSet.Error())
	default:
		http.InternalServerError(c)
	}
}

func (ui *UI) DeleteCategory(c *http.Context) {
	err := ui.app.RemoveCategoryByID(c.Request.Path.Params("category"))
	switch err {
	case nil:
		http.NoContent(c)
	case service.ErrNoSuchCategory:
		c.String(http.StatusNotFound, service.ErrNoSuchCategory.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}
