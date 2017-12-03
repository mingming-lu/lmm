package main

import (
	"./articles"
	"./image"
	"./profile"

	"github.com/akinaru-lu/elesion"
)

func main() {
	el := elesion.Default()
	el.Handle("/articles", articles.Handler)
	el.Handle("/photos", image.Handler)
	el.Handle("/profile", profile.Handler)
	el.Run(":8081")
}
