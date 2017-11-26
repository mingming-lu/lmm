package main

import (
    "log"
    "net/http"
    "os"
)

const Res = "image/res"
const Special = "image/special"

func handleImage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, Res + r.URL.Path)
}

func avatar(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, Special + "/avatar.jpg")
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
    http.HandleFunc("/", handleImage)
    http.HandleFunc("/avatar", avatar)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
