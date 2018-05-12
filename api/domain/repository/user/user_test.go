package user

import (
	"lmm/api/db"
	model "lmm/api/domain/model/user"
	"lmm/api/testing"
)

func TestSave(t *testing.T) {
	tester := testing.NewTester(t)

	repo := New()
	testing.InitTable("user")
	m := model.New("foobar", "1234")
	user, err := repo.Save(m)
	tester.NoError(err)
	tester.Is(user.ID, uint64(1))
	tester.Is(user.Name, "foobar")

	db := db.Default()
	defer db.Close()

	stmt := db.Must("SELECT * FROM user WHERE id = ?")
	defer stmt.Close()

	r := stmt.QueryRow(user.ID)
	r.Scan(&m)
	tester.Is(m.ID, uint64(1))
}
