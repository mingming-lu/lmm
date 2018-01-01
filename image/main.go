package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akinaru-lu/elesion"
)

const Res = "image/res"
const Special = "image/special"

func handler(c *elesion.Context) {
	c.File(Res + c.Request.URL.Path)
}

func avatar(c *elesion.Context) {
	c.File(Special + "/avatar.jpg")
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func init() {
	ensureDir(Res)
	ensureDir(Special)
}

func main() {
	router := elesion.Default("[image]")
	router.GET("/", handler)
	router.GET("/avatar", avatar)
	log.Fatal(http.ListenAndServe(":8082", router))
}
