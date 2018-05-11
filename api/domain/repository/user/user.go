package user

import (
	"lmm/api/db"
	model "lmm/api/domain/model/user"
)

type Repository struct {
	db func() *db.DB
}

func New() *Repository {
	return &Repository{db: func() *db.DB { return db.Default() }}
}

// Save return a User model with generated id
func (repo *Repository) Save(user *model.User) (*model.User, error) {
	db := repo.db()
	defer db.Close()

	stmt := db.Must(`INSERT INTO user (name, password, guid, token, created_at) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Password, user.GUID, user.Token, user.CreatedAt.UTC())
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = uint64(id)
	return user, nil
}

func Add(name, password, guid, token string) (int64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must(`INSERT INTO user (name, password, guid, token) values (?, ?, ?, ?)`)
	defer stmt.Close()

	res, err := stmt.Exec(name, password, guid, token)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func ByName(name string) (*model.User, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, name, password, guid, token, created_at FROM user WHERE name = ?")
	defer stmt.Close()

	user := model.User{}
	err := stmt.QueryRow(name).Scan(
		&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt,
	)

	return &user, err
}

func ByToken(token string) (*model.User, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, name, password, guid, token, created_at FROM user WHERE token = ?")
	defer stmt.Close()

	user := model.User{}
	err := stmt.QueryRow(token).Scan(
		&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt,
	)

	return &user, err
}
