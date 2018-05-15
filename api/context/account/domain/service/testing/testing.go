package testing

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/db"
	"lmm/api/utils/uuid"
)

func NewUser() *model.User {
	db := db.Default()
	defer db.Close()

	stmt1 := db.Must(`INSERT INTO user (name, password, guid, token, created_at) VALUES(?, ?, ?, ?, ?)`)
	defer stmt1.Close()

	name := uuid.New()[:32]
	password := uuid.New()
	user := model.New(name, password)

	result, err := stmt1.Exec(user.Name, user.Password, user.GUID, user.Token, user.CreatedAt)
	if err != nil {
		panic(err)
	}
	userID, err := result.LastInsertId()

	stmt2 := db.Must(`SELECT id, name, guid, token, created_at FROM user WHERE id = ?`)
	defer stmt2.Close()

	user = &model.User{}
	if err := stmt2.QueryRow(userID).Scan(&user.ID, &user.Name, &user.GUID, &user.Token, &user.CreatedAt); err != nil {
		panic(err)
	}
	user.Password = password
	return user
}
