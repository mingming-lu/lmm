package main

import (
	"lmm/api"
	"lmm/api/http"
	"lmm/api/storage"
)

func main() {
	db := storage.NewDB()
	router := api.NewRouter(db, nil)

	http.Serve(":8002", router)
}
