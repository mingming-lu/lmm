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

func index(c *elesion.Context) {
	// TODO handle query parameter `?thumbnail=true`
	c.File(dirRaw + c.Request.URL.Path)
}

func ensureDir(name string) string {
	path := where() + "/" + name

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	} else if err != nil {
		log.Fatal(err)
	}

	return path
}

var (
	dirRaw       string
	dirThumbnail string
)

func init() {
	dirRaw = ensureDir("raw")
	dirThumbnail = ensureDir("thumbnail")
}

func main() {
	router := elesion.Default("[image]")
	router.GET("/:name", index)
	log.Fatal(http.ListenAndServe(":8003", router))
}
