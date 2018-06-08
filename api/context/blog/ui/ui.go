package ui

import (
	"fmt"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
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

	app := appservice.New(repository.NewBlogRepository())
	blogID, err := app.PostNewBlog(user.ID(), blog.Title, blog.Text)
	if err != nil {
		http.InternalServerError(c)
		return
	}

	c.Header("Location", fmt.Sprintf("/blog/%d", blogID))
	c.Status(http.StatusCreated).String("success")
}
