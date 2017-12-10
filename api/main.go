package main

import (
	"./articles"
	"./image"
	"./profile"

	"github.com/akinaru-lu/elesion"
)

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
