package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Images     []Image `json:"images"`
	NextCursor string  `json:"next_cursor"`
}

type Image struct {
	URL string `json:"url"`
}

const Host = "http://localhost:8082/"
const Path = "image/res/"

func HandlePhotos(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		fmt.Println(err)
	}

	var photos []Image
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		u := Host + file.Name()
		photos = append(photos, Image{URL: u})
	}

	resp := Response{
		Images: photos,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
