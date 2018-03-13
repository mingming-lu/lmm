package image

import (
	"lmm/api/db"
	model "lmm/api/domain/model/image"
)

func FetchAllPhotos(userID int64) ([]model.Minimal, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT i.name FROM photo AS p INNER JOIN image AS i ON p.image = i.id AND p.user = i.user WHERE p.user = ? ORDER BY created_at DESC")
	defer stmt.Close()

	images := make([]model.Minimal, 0)

	itr, err := stmt.Query(userID)
	if err != nil {
		return images, err
	}

	for itr.Next() {
		image := model.Minimal{}
		err = itr.Scan(&image.Name)
		if err != nil {
			return images, err
		}
		images = append(images, image)
	}

	return images, nil
}

func MarkAsPhoto(userID, imageID int64) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO photo (user, image) VALUES (?, ?)")
	defer stmt.Close()

	_, err := stmt.Exec(userID, imageID)
	if err != nil {
		return err
	}

	return nil
}
