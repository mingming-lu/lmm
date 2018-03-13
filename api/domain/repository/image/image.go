package image

import (
	"io/ioutil"
	"lmm/api/db"
	model "lmm/api/domain/model/image"
	"os"
)

const pathRaw = "image/raw/"

func Add(userID int64, data []model.ImageData) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO image (user, name) VALUES (?, ?)")
	defer stmt.Close()

	tx, err := d.Begin()
	if err != nil {
		return err
	}
	stmt = tx.Stmt(stmt)

	for _, image := range data {
		_, err = stmt.Exec(userID, image.Name)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = ioutil.WriteFile(pathRaw+image.Name, image.Data, os.ModePerm)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func FetchAllImage(userID int64) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT name FROM image WHERE user = ? ORDER BY created_at DESC")
	defer stmt.Close()

	itr, err := stmt.Query(userID)
	if err != nil {
		return
	}
	for itr.Next() {
	}
}
