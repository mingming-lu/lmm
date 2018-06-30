package infra

import (
	"fmt"
	"lmm/api/context/blog/domain/model"
	"lmm/api/storage"
)

type TagStorage struct {
	db *storage.DB
}

func NewTagStorage(db *storage.DB) *TagStorage {
	return &TagStorage{db: db}
}

func (s *TagStorage) Add(tag *model.Tag) error {
	stmt := s.db.MustPrepare(`INSERT INTO tag (id, blog, name) VALUES(?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(tag.ID(), tag.BlogID(), tag.Name())
	return err
}

func (s *TagStorage) FindByID(id uint64) (*model.Tag, error) {
	return s.selectRow(`SELECT id, blog, name FROM tag WHERE id = ?`, id)
}

func (s *TagStorage) FindAll() ([]*model.Tag, error) {
	return s.selectRows(`SELECT id, blog, name FROM tag`)
}

func (s *TagStorage) FindAllByBlog(blog *model.Blog) ([]*model.Tag, error) {
	return s.selectRows(`SELECT id, blog, name FROM tag WHERE blog = ?`, blog.ID())
}

func (s *TagStorage) Update(tag *model.Tag) error {
	return s.updateRow(`UPDATE FROM tag SET name = ? WHERE id = ?`, tag.Name(), tag.ID())
}

func (s *TagStorage) Remove(tag *model.Tag) error {
	return s.updateRow(`DELETE FROM tag WHERE id = ?`, tag.ID())
}

func (s *TagStorage) selectRow(query string, args ...interface{}) (*model.Tag, error) {
	stmt := s.db.MustPrepare(query)
	defer stmt.Close()

	var (
		blogID  uint64
		tagID   uint64
		tagName string
	)

	err := stmt.QueryRow(args...).Scan(&tagID, &blogID, &tagName)
	if err != nil {
		return nil, err
	}

	return model.NewTag(tagID, blogID, tagName)
}

func (s *TagStorage) selectRows(query string, args ...interface{}) ([]*model.Tag, error) {
	stmt := s.db.MustPrepare(query)
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		blogID  uint64
		tagID   uint64
		tagName string
	)

	tags := make([]*model.Tag, 0)
	for rows.Next() {
		err := rows.Scan(&tagID, &tagName)
		if err != nil {
			return nil, err
		}
		tag, err := model.NewTag(tagID, blogID, tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *TagStorage) updateRow(query string, args ...interface{}) error {
	stmt := s.db.MustPrepare(query)
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("rows affected is expected to be 1 but got %d", rowsAffected)
	}
	return nil
}
