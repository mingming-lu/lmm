package main

import (
	"lmm/api"
	"lmm/api/http"
	"lmm/api/storage"
)

func main() {
	db := storage.NewDB()
	cache := storage.NewCacheEngine()
	router := api.NewRouter(db, cache)

	http.Serve(":8002", router)
}
