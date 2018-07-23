package ui

import (
	"fmt"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
)

type UI struct {
	app *appservice.AppService
}

func New(
	blogRepo repository.BlogRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) *UI {
	app := appservice.New(blogRepo, categoryRepo, tagRepo)
	return &UI{app: app}
}

func (ui *UI) PostBlog(c *http.Context) {
	user := c.Values().Get("user").(*account.User)
	blogID, err := ui.app.PostNewBlog(user, c.Request.Body)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/blog/%d", blogID)).String(http.StatusCreated, "success")
	case domain.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, domain.ErrEmptyBlogTitle.Error())
	default:
		c.Logger().Error(err.Error())
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
	case domain.ErrInvalidCount, domain.ErrInvalidPage:
		c.String(http.StatusBadRequest, err.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) GetBlog(c *http.Context) {
	blog, err := ui.app.GetBlogByID(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, blog)
	case domain.ErrNoSuchBlog:
		c.String(http.StatusNotFound, domain.ErrNoSuchBlog.Error())
	default:
		c.Logger().Error(err.Error())
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
	case domain.ErrBlogNoChange:
		http.NoContent(c)
	case domain.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, domain.ErrEmptyBlogTitle.Error())
	case domain.ErrNoPermission:
		c.String(http.StatusForbidden, domain.ErrNoSuchBlog.Error())
	case domain.ErrNoSuchBlog:
		c.String(http.StatusNotFound, domain.ErrNoSuchBlog.Error())
	default:
		c.Logger().Error(err.Error())
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
	case domain.ErrNoSuchBlog, domain.ErrNoSuchCategory:
		c.String(http.StatusBadRequest, err.Error())
	default:
		c.Logger().Error(err.Error())
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
	case domain.ErrInvalidCategoryName, domain.ErrDuplicateCategoryName:
		c.String(http.StatusBadRequest, err.Error())
	default:
		c.Logger().Error(err.Error())
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
	case domain.ErrCategoryNoChanged:
		http.NoContent(c)
	case domain.ErrInvalidCategoryName:
		c.String(http.StatusBadRequest, domain.ErrInvalidCategoryName.Error())
	case domain.ErrNoSuchCategory:
		c.String(http.StatusNotFound, domain.ErrNoSuchCategory.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) GetAllCategoris(c *http.Context) {
	categories, err := ui.app.GetAllCategories()
	switch err {
	case nil:
		c.JSON(http.StatusOK, categories)
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) GetBlogCagetory(c *http.Context) {
	category, err := ui.app.GetCategoryOfBlog(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, category)
	case domain.ErrNoSuchBlog:
		c.String(http.StatusNotFound, domain.ErrNoSuchBlog.Error())
	case domain.ErrCategoryNotSet:
		c.String(http.StatusNotFound, domain.ErrCategoryNotSet.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) DeleteCategory(c *http.Context) {
	err := ui.app.RemoveCategoryByID(c.Request.Path.Params("category"))
	switch err {
	case nil:
		http.NoContent(c)
	case domain.ErrNoSuchCategory:
		c.String(http.StatusNotFound, domain.ErrNoSuchCategory.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) NewBlogTag(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	tag := Tag{}
	if err := c.Request.ScanBody(&tag); err != nil {
		c.Logger().Warn(err.Error())
		http.BadRequest(c)
		return
	}

	err := ui.app.AddNewTagToBlog(user, c.Request.Path.Params("blog"), tag.Name)
	switch err {
	case nil:
		c.String(http.StatusCreated, "success")
	case domain.ErrInvalidTagName:
		c.String(http.StatusBadRequest, domain.ErrInvalidTagName.Error())
	case domain.ErrNoSuchBlog:
		c.String(http.StatusNotFound, domain.ErrNoSuchBlog.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) UpdateTag(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	tag := Tag{}
	if err := c.Request.ScanBody(&tag); err != nil {
		c.Logger().Warn(err.Error())
		http.BadRequest(c)
		return
	}

	err := ui.app.UpdateBlogTag(user, c.Request.Path.Params("tag"), tag.Name)
	switch err {
	case nil:
		c.String(http.StatusOK, "success")
	case domain.ErrInvalidTagName:
		c.String(http.StatusBadRequest, domain.ErrInvalidTagName.Error())
	case domain.ErrNoSuchTag:
		c.String(http.StatusNotFound, domain.ErrNoSuchTag.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) DeleteTag(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	err := ui.app.RemoveBlogTag(user, c.Request.Path.Params("tag"))
	switch err {
	case nil:
		http.NoContent(c)
	case domain.ErrNoSuchTag:
		c.String(http.StatusOK, domain.ErrNoSuchTag.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) GetAllTags(c *http.Context) {
	tags, err := ui.app.GetAllTags()
	switch err {
	case nil:
		c.JSON(http.StatusOK, tagsToJSON(tags))
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) GetAllTagsOfBlog(c *http.Context) {
	tags, err := ui.app.GetAllTagsOfBlog(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, tagsToJSON(tags))
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}
