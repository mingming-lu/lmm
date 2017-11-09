package main

import (
    "log"
    "net/http"
)

const ImagePath = "image/res"

func handleImage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, ImagePath + r.URL.Path)
}

func main() {
    http.HandleFunc("/", handleImage)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
