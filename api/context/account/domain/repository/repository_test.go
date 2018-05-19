package repository

import (
	"lmm/api/context/account/domain/model"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/db"
	"lmm/api/testing"
	"lmm/api/utils/sha256"
)

func TestSave(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	repo := New()
	m := model.NewUser("foobar", "1234")
	user, err := repo.Put(m)
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

func TestFindByName_Success(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	user := testingService.NewUser()
	user.Password = sha256.Hex([]byte(user.GUID + user.Password))

	repo := New()
	found, err := repo.FindByName(user.Name)

	tester.NoError(err)
	tester.Isa(&model.User{}, found)
	tester.Is(user, found)
}

func TestFindByName_NotFound(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	repo := New()
	found, err := repo.FindByName("foo")

	tester.Error(err)
	tester.Nil(found)
	tester.Is(db.ErrNoRows, err)
}

func TestFindByToken_Success(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	user := testingService.NewUser()
	repo := New()

	found, err := repo.FindByToken(user.Token)
	tester.NoError(err)
	tester.Isa(&model.User{}, found)
	tester.Is(user.Name, found.Name)
}

func TestFindByToken_NotFound(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	repo := New()

	found, err := repo.FindByToken("1234")
	tester.Error(err)
	tester.Is(db.ErrNoRows, err)
	tester.Nil(found)
}
