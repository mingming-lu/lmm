package ui

import (
	"fmt"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
	"lmm/api/utils/strings"
	"log"
)

func PostBlog(c *http.Context) {
	blog := Blog{}
	err := c.Request.ScanBody(&blog)
	if err != nil {
		http.BadRequest(c)
		log.Println(err)
		return
	}

	user := c.Values().Get("user").(*account.User)

	app := appservice.NewBlogApp(repository.NewBlogRepository())
	blogID, err := app.PostNewBlog(user.ID(), blog.Title, blog.Text)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/blog/%d", blogID)).String(http.StatusCreated, "success")
	case appservice.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, appservice.ErrEmptyBlogTitle.Error())
	default:
		http.InternalServerError(c)
	}
}

func GetAllBlog(c *http.Context) {
	app := appservice.NewBlogApp(repository.NewBlogRepository())
	blogItems, page, hasNextPage, err := app.FindAllBlog(
		c.Request.Query("count"),
		c.Request.Query("page"),
	)

	switch err {
	case nil:
		blogList := make([]BlogResponse, len(blogItems))
		for index, blogItem := range blogItems {
			blogList[index].ID = strings.Uint64ToStr(blogItem.ID())
			blogList[index].Title = blogItem.Title()
			blogList[index].Text = blogItem.Text()
			blogList[index].CreatedAt = blogItem.CreatedAt().UTC().String()
			blogList[index].UpdatedAt = blogItem.UpdatedAt().UTC().String()
		}
		c.JSON(http.StatusOK, BlogListResponse{
			Blog:        blogList,
			Page:        page,
			HasNextPage: hasNextPage,
		})
	case appservice.ErrInvalidCount, appservice.ErrInvalidPage:
		c.String(http.StatusBadRequest, err.Error())
	default:
		http.InternalServerError(c)
	}
}

func GetBlog(c *http.Context) {
	app := appservice.NewBlogApp(repository.NewBlogRepository())
	blog, err := app.FindBlogByID(c.Request.Path.Params("blog"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, BlogResponse{
			ID:        strings.Uint64ToStr(blog.ID()),
			Title:     blog.Title(),
			Text:      blog.Text(),
			CreatedAt: blog.CreatedAt().UTC().String(),
			UpdatedAt: blog.UpdatedAt().UTC().String(),
		})
	case appservice.ErrNoSuchBlog:
		c.String(http.StatusNotFound, appservice.ErrNoSuchBlog.Error())
	default:
		http.InternalServerError(c)
	}
}

func UpdateBlog(c *http.Context) {
	user := c.Values().Get("user").(*account.User)
	app := appservice.NewBlogApp(repository.NewBlogRepository())

	blog := Blog{}
	c.Request.ScanBody(&blog)

	err := app.EditBlog(user.ID(), c.Request.Path.Params("blog"), blog.Title, blog.Text)
	switch err {
	case nil:
		c.String(http.StatusOK, "success")
	case appservice.ErrBlogNoChange:
		http.NoContent(c)
	case appservice.ErrEmptyBlogTitle:
		c.String(http.StatusBadRequest, appservice.ErrEmptyBlogTitle.Error())
	case appservice.ErrNoPermission:
		c.String(http.StatusForbidden, appservice.ErrNoSuchBlog.Error())
	case appservice.ErrNoSuchBlog:
		c.String(http.StatusNotFound, appservice.ErrNoSuchBlog.Error())
	default:
		http.InternalServerError(c)
	}
}
