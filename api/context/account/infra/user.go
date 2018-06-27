package infra

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/storage"
	"time"
)

type UserStorage struct {
	db *storage.DB
}

func NewUserStorage(db *storage.DB) *UserStorage {
	return &UserStorage{db: db}
}

// Put puts a new user intoUserStorage)itory and return a User model with generated id
func (s *UserStorage) Add(user *model.User) error {
	stmt := s.db.MustPrepare(`INSERT INTO user (id, name, password, token, created_at) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(user.ID(), user.Name(), user.Password(), user.Token(), user.CreatedAt().UTC())
	if err != nil {
		return err
	}
	return nil
}

// FindByName return a user model determined by name
func (s *UserStorage) FindByName(name string) (*model.User, error) {
	stmt := s.db.MustPrepare(`SELECT id, name, password, token, created_at FROM user WHERE name = ?`)
	defer stmt.Close()

	var (
		userID        uint64
		userName      string
		userPassword  string
		userToken     string
		userCreatedAt time.Time
	)
	err := stmt.QueryRow(name).Scan(&userID, &userName, &userPassword, &userToken, &userCreatedAt)
	if err != nil {
		return nil, err
	}
	return model.NewUser(userID, userName, userPassword, userToken, userCreatedAt), nil
}

func (s *UserStorage) FindByToken(token string) (*model.User, error) {
	stmt := s.db.MustPrepare(`SELECT id, name, password, token, created_at FROM user WHERE token = ?`)
	defer stmt.Close()

	var (
		userID        uint64
		userName      string
		userPassword  string
		userToken     string
		userCreatedAt time.Time
	)
	err := stmt.QueryRow(token).Scan(&userID, &userName, &userPassword, &userToken, &userCreatedAt)
	if err != nil {
		return nil, err
	}
	return model.NewUser(userID, userName, userPassword, userToken, userCreatedAt), nil
}
