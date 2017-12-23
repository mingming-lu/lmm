package main

import (
	"lmm/api/articles"
	"lmm/api/image"
	"lmm/api/profile"

	"github.com/akinaru-lu/elesion"
	"lmm/api/db"
)

func init() {
	db.New().Create("lmm").Close()
}

func Allowed(c *elesion.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	el := elesion.Default()
	el.Use(Allowed)
	el.Handle("/articles", articles.Handler)
	el.Handle("/photos", image.Handler)
	el.Handle("/profile", profile.Handler)
	el.Run(":8081")
}
