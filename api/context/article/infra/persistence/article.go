package persistence

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lmm/api/context/article/domain/model"
	"lmm/api/storage"
)

// ArticleStorage is a implementation of ArticleRepository
type ArticleStorage struct {
	db *storage.DB
}

// NewArticleStorage constructs a new article repository with concrete struct
func NewArticleStorage(db *storage.DB) *ArticleStorage {
	return &ArticleStorage{db: db}
}

// NextID generate a random string
func (s *ArticleStorage) NextID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.New().String())))[:8]
}

// Save persist a article domain model
func (s *ArticleStorage) Save(article *model.Article) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	saveArticle, err := tx.Prepare(`
		insert into article (uid, user, title, body, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?)
		on duplicate key update id = LAST_INSERT_ID(id), title = ?, body = ?, updated_at = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	deleteTags, err := tx.Prepare(`
		delete at from article_tag at left join article a on a.id = at.article_id where at.article_id = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	saveTags, err := tx.Prepare("" +
		"insert into article_tag (article_id, `order`, name) values (?, ?, ?)" +
		"")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	now := time.Now()

	if err != nil {
		return nil
	}

	res, err := saveArticle.Exec(
		article.ID().String(),
		article.Author().ID(),
		article.Content().Text().Title(),
		article.Content().Text().Body(),
		now,
		now,
		article.Content().Text().Title(),
		article.Content().Text().Body(),
		now,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if _, err := deleteTags.Exec(lastID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	for _, tag := range article.Content().Tags() {
		if _, err := saveTags.Exec(lastID, tag.ID().Order(), tag.Name()); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

// Remove is not implemented
func (s *ArticleStorage) Remove(article *model.Article) error {
	panic("not implemented")
}

// FindByID returns a article domain model by given id if exists
func (s *ArticleStorage) FindByID(id *model.ArticleID) (*model.Article, error) {
	stmt := s.db.MustPrepare("SELECT uid, writer, title, text FROM article WHERE uid = ?")
	defer stmt.Close()

	return s.modelFromRow(stmt.QueryRow(id.String()))
}

func (s *ArticleStorage) modelFromRow(row *sql.Row) (*model.Article, error) {
	panic("not implemented")
}
