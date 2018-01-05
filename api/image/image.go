package image

import (
	"io/ioutil"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

type Response struct {
	Images     []Image `json:"images"`
	NextCursor string  `json:"next_cursor"`
}

type Image struct {
	URL string `json:"url"`
}

const Host = "http://image.lmm.local/photos/"
const Path = "image/photos/"

func Handler(c *elesion.Context) {
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var photos []Image
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		u := Host + file.Name()
		photos = append(photos, Image{URL: u})
	}

	data := Response{
		Images: photos,
	}
	c.Status(http.StatusOK).JSON(data)
}
