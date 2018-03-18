package image

import (
	"lmm/api/db"
	model "lmm/api/domain/model/image"
)

func FetchAllPhotos(userID int64) ([]model.Minimal, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must(`
		SELECT i.name
		FROM photo AS p
		INNER JOIN image AS i ON p.image = i.id AND p.user = i.user
		WHERE p.user = ? AND p.deleted = 0
		ORDER BY p.last_modified DESC
	`)
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

func SavePhoto(userID, imageID int64, shown bool) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must(`
		INSERT INTO photo (user, image)
		VALUES (?, ?)
		ON DUPLICATE key update deleted = ?
	`)
	defer stmt.Close()

	var err error
	if shown {
		_, err = stmt.Exec(userID, imageID, 0)
	} else {
		_, err = stmt.Exec(userID, imageID, 1)
	}
	if err != nil {
		return err
	}

	return nil
}
