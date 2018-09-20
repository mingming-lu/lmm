package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/service"
	"lmm/api/storage"
	"lmm/api/utils/strings"
)

// ArticleStorage is a implementation of ArticleRepository
type ArticleStorage struct {
	db            *storage.DB
	authorService service.AuthorService
}

// NewArticleStorage constructs a new article repository with concrete struct
func NewArticleStorage(db *storage.DB, authorService service.AuthorService) *ArticleStorage {
	return &ArticleStorage{db: db, authorService: authorService}
}

// NextID generate a random string
func (s *ArticleStorage) NextID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Save persist a article domain model
func (s *ArticleStorage) Save(article *model.Article) error {
	stmt := s.db.MustPrepare("insert into article (uid, user, title, body, created_at, updated_at) " +
		"values (?, ?, ?, ?, ?, ?) " +
		"on duplicate key update title = ?, body = ?, updated_at = ?",
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
	stmt := s.db.MustPrepare("SELECT uid, user, title, body FROM article WHERE uid = ?")
	defer stmt.Close()

	return s.userModelFromRow(stmt.QueryRow(id.String()))
}

func (s *ArticleStorage) userModelFromRow(row *sql.Row) (*model.Article, error) {
	var (
		rawArticleID string
		userID       uint64
		title        string
		body         string
	)
	if err := row.Scan(&rawArticleID, &userID, &title, &body); err != nil {
		return nil, err
	}
	author, err := s.authorService.AuthorFromUserID(userID)
	if err != nil {
		return nil, err
	}
	articleID, err := model.NewArticleID(rawArticleID)
	if err != nil {
		return nil, err
	}
	text, err := model.NewText(title, body)
	if err != nil {
		return nil, err
	}
	return model.NewArticle(articleID, text, author, make([]*model.Tag, 0)), nil
}
