package infra

import (
	"database/sql"
	"lmm/api/context/article/domain/model"
	"lmm/api/storage"
	"lmm/api/utils/strings"
	"time"

	"github.com/google/uuid"
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
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Save persist a article domain model
func (s *ArticleStorage) Save(article *model.Article) error {
	stmt := s.db.MustPrepare("INSERT INTO article uid, writer, title, text, created_at, updated_at" +
		"VALUES (?, ?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE SET title = ?, SET text = ? updated_at = ?",
	)
	defer stmt.Close()

	now := time.Now()

	_, err := stmt.Exec(
		article.ID().String(),
		article.Author().ID(),
		article.Text().Title(),
		article.Text().Body(),
		now,
		now,
		article.Text().Title(),
		article.Text().Body(),
		now,
	)

	return err
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
