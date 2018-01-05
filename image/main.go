package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/akinaru-lu/elesion"
)

func where() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get caller infomation")
	}
	return path.Dir(filename)
}

func photos(c *elesion.Context) {
	c.File(here + c.Request.URL.Path)
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

var here = "image"

func init() {
	here := where()
	ensureDir(here + "/photos")
	ensureDir(here + "/special")
}

func main() {
	router := elesion.Default("[image]")
	router.GET("/photos/:name", photos)
	log.Fatal(http.ListenAndServe(":8082", router))
}
