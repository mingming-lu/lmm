package main

import (
    "log"
    "net/http"
    "os"
)

const ImageResPath = "image/res"

func handleImage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, ImageResPath + r.URL.Path)
}

func ensureRes() {
    if _, err := os.Stat(ImageResPath); os.IsNotExist(err) {
        os.Mkdir(ImageResPath, os.ModePerm)
    }
}

func init() {
    ensureRes()
}

func main() {
    http.HandleFunc("/", handleImage)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
