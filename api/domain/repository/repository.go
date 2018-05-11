package repository

import "lmm/api/db"

type Repository struct {
	db func() *db.DB
}

func New() *Repository {
	return &Repository{db: func() *db.DB { return db.Default() }}
}
