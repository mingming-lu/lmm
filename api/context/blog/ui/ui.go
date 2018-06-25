package ui

import (
	"fmt"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/appservice"
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
	case appservice.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, appservice.ErrEmptyBlogTitle.Error())
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
	case appservice.ErrInvalidCount, appservice.ErrInvalidPage:
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
	case appservice.ErrNoSuchBlog:
		c.String(http.StatusNotFound, appservice.ErrNoSuchBlog.Error())
	default:
		log.Println(err)
		http.InternalServerError(c)
	}
}

// func UpdateBlog(c *http.Context) {
// 	user, ok := c.Values().Get("user").(*account.User)
// 	if !ok {
// 		http.Unauthorized(c)
// 		return
// 	}
// 	app := appservice.NewBlogApp(repository.NewBlogRepository())
//
// 	blog := Blog{}
// 	err := c.Request.ScanBody(&blog)
// 	if err != nil {
// 		log.Println(err)
// 		http.BadRequest(c)
// 		return
// 	}
//
// 	err = app.EditBlog(user.ID(), c.Request.Path.Params("blog"), blog.Title, blog.Text)
// 	switch err {
// 	case nil:
// 		c.String(http.StatusOK, "success")
// 	case appservice.ErrBlogNoChange:
// 		http.NoContent(c)
// 	case appservice.ErrEmptyBlogTitle:
// 		c.String(http.StatusBadRequest, appservice.ErrEmptyBlogTitle.Error())
// 	case appservice.ErrNoPermission:
// 		c.String(http.StatusForbidden, appservice.ErrNoSuchBlog.Error())
// 	case appservice.ErrNoSuchBlog:
// 		c.String(http.StatusNotFound, appservice.ErrNoSuchBlog.Error())
// 	default:
// 		log.Println(err)
// 		http.InternalServerError(c)
// 	}
// }
//
// func SetBlogCategory(c *http.Context) {
// 	_, ok := c.Values().Get("user").(*account.User)
// 	if !ok {
// 		http.Unauthorized(c)
// 		return
// 	}
//
// 	category := Category{}
// 	if err := c.Request.ScanBody(&category); err != nil {
// 		log.Println(err)
// 		http.BadRequest(c)
// 		return
// 	}
//
// 	err := app.SetBlogCategory(c.Request.Path.Params("blog"), category.Name)
// 	switch err {
// 	case nil:
// 		c.String(http.StatusOK, "success")
// 	case service.ErrNoSuchBlog, service.ErrNoSuchCategory:
// 		c.String(http.StatusBadRequest, err.Error())
// 	default:
// 		http.InternalServerError(c)
// 	}
// }
//
// func PostCategory(c *http.Context) {
// 	app := appservice.NewCategoryApp(repository.NewCategoryRepository())
//
// 	category := Category{}
// 	err := c.Request.ScanBody(&category)
// 	if err != nil {
// 		log.Println(err)
// 		http.BadRequest(c)
// 		return
// 	}
//
// 	categoryID, err := app.AddNewCategory(category.Name)
// 	switch err {
// 	case nil:
// 		c.Header("Location", fmt.Sprintf("/categories/%d", categoryID)).String(http.StatusCreated, "success")
// 	case appservice.ErrInvalidCategoryName, appservice.ErrDuplicateCategoryName:
// 		c.String(http.StatusBadRequest, err.Error())
// 	default:
// 		log.Println(err)
// 		http.InternalServerError(c)
// 	}
// }
//
// func UpdateCategory(c *http.Context) {
// 	app := appservice.NewCategoryApp(repository.NewCategoryRepository())
//
// 	categoryID := c.Request.Path.Params("category")
// 	category := Category{}
// 	err := c.Request.ScanBody(&category)
// 	if err != nil {
// 		log.Println(err)
// 		http.BadRequest(c)
// 		return
// 	}
//
// 	err = app.UpdateCategoryName(categoryID, category.Name)
//
// 	switch err {
// 	case nil:
// 		c.String(http.StatusOK, "success")
// 	case appservice.ErrCategoryNoChanged:
// 		http.NoContent(c)
// 	case appservice.ErrInvalidCategoryName:
// 		c.String(http.StatusBadRequest, appservice.ErrInvalidCategoryName.Error())
// 	case appservice.ErrNoSuchCategory:
// 		c.String(http.StatusNotFound, appservice.ErrNoSuchCategory.Error())
// 	default:
// 		log.Println(err)
// 		http.InternalServerError(c)
// 	}
// }
//
// func GetAllCategoris(c *http.Context) {
// 	app := appservice.NewCategoryApp(repository.NewCategoryRepository())
//
// 	models, err := app.FindAllCategories()
// 	switch err {
// 	case nil:
// 		categories := make([]*Category, len(models))
// 		for index, model := range models {
// 			categories[index].Name = model.Name()
// 		}
// 		c.JSON(http.StatusOK, CategoriesResponse{
// 			Categories: categories,
// 		})
// 	default:
// 		log.Println(err)
// 		http.InternalServerError(c)
// 	}
// }
//
// func GetBlogCagetory(c *http.Context) {
// 	category, err := app.GetCategoryOfBlog(c.Request.Path.Params("blog"))
// 	switch err {
// 	case nil:
// 		c.JSON(http.StatusOK, category)
// 	case service.ErrNoSuchBlog:
// 		c.String(http.StatusNotFound, service.ErrNoSuchBlog.Error())
// 	case service.ErrCategoryNotSet:
// 		c.String(http.StatusNotFound, service.ErrCategoryNotSet.Error())
// 	default:
// 		http.InternalServerError(c)
// 	}
// }
//
// func DeleteCategory(c *http.Context) {
// 	app := appservice.NewCategoryApp(repository.NewCategoryRepository())
//
// 	categoryID := c.Request.Path.Params("category")
//
// 	err := app.Remove(categoryID)
// 	switch err {
// 	case nil:
// 		http.NoContent(c)
// 	case appservice.ErrNoSuchCategory:
// 		c.String(http.StatusNotFound, appservice.ErrNoSuchCategory.Error())
// 	default:
// 		log.Println(err)
// 		http.InternalServerError(c)
// 	}
// }
