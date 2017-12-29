package main

import (
	"lmm/api/articles"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"

	"github.com/akinaru-lu/elesion"
)

func init() {
	db.Init()
}

func Allowed(c *elesion.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	el := elesion.Default()
	el.Use(Allowed)
	el.Handle("/articles", articles.GetArticles)
	el.Handle("/article", articles.GetArticle)
	el.Handle("/photos", image.Handler)
	el.Handle("/profile", profile.Handler)
	el.Run(":8081")
}
