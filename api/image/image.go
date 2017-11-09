package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Result     Result `json:"result"`
}

type Result struct {
	Images []Image `json:"images"`
}

type Image struct {
	URL string `json:"url"`
}

const ImageHost = "http://localhost:8082/"
const ImagePath = "image/res/"

func HandlePhotos(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(ImagePath)
	if err != nil {
		fmt.Println(err)
	}

	var photos []Image
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		u := ImageHost + file.Name()
		photos = append(photos, Image{URL: u})
	}

	resp := Response{
		StatusCode: http.StatusOK,
		Result: Result{
			Images: photos,
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
