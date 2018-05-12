package user

import (
	"lmm/api/db"
	model "lmm/api/domain/model/user"
	"lmm/api/testing"
)

func TestSave(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	repo := New()
	m := model.New("foobar", "1234")
	user, err := repo.Save(m)
	tester.NoError(err)
	tester.Is(uint64(1), user.ID)
	tester.Is("foobar", user.Name)

	db := db.Default()
	defer db.Close()

	stmt := db.Must("SELECT * FROM user WHERE id = ?")
	defer stmt.Close()

	r := stmt.QueryRow(user.ID)
	r.Scan(&m)
	tester.Is(uint64(1), m.ID)
}
