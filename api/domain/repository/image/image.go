package image

import (
	"io/ioutil"
	"lmm/api/db"
	model "lmm/api/domain/model/image"
	"os"
)

const pathRaw = "image/raw/"

func Add(userID int64, imageType model.ImageType, name string, data []byte) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO image (user, type, url) VALUES (?, ?, ?)")
	defer stmt.Close()

	tx, err := d.Begin()
	if err != nil {
		return err
	}
	stmt = tx.Stmt(stmt)

	_, err = stmt.Exec(userID, imageType, model.BaseURL+name)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ioutil.WriteFile(pathRaw+name, data, os.ModePerm)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func Fetch(userID int64, t model.ImageType) ([]model.Image, error) {
	return nil, nil
}
