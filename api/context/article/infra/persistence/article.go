package persistence

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/service"
	"lmm/api/storage"
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
func (s *ArticleStorage) NextID(c context.Context) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.New().String())))[:8]
}

// Save persist a article domain model
func (s *ArticleStorage) Save(c context.Context, article *model.Article) error {
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
		delete at from article_tag at left join article a on a.id = at.article where at.article = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	saveTags, err := tx.Prepare(`
		insert into article_tag (article, sort, name) values (?, ?, ?)
	`)
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
func (s *ArticleStorage) Remove(c context.Context, article *model.Article) error {
	panic("not implemented")
}

// FindByID returns a article domain model by given id if exists
func (s *ArticleStorage) FindByID(c context.Context, id *model.ArticleID) (*model.Article, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`select id, uid, user, title, body from article where uid = ? for update`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	article, err := s.userModelFromRow(context.TODO(), stmt.QueryRow(id.String()))
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleStorage) userModelFromRow(c context.Context, row *sql.Row) (*model.Article, error) {
	var (
		id           int
		rawArticleID string
		userID       uint64
		title        string
		body         string
	)
	if err := row.Scan(&id, &rawArticleID, &userID, &title, &body); err != nil {
		if err == storage.ErrNoRows {
			return nil, domain.ErrNoSuchArticle
		}
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

	stmt := s.db.MustPrepare("select sort, name from article_tag where article = ?")
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		order uint
		name  string
	)

	tags := make([]*model.Tag, 0)
	for rows.Next() {
		if err := rows.Scan(&order, &name); err != nil {
			return nil, err
		}
		tag, err := model.NewTag(articleID, order, name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	content, err := model.NewContent(text, tags)
	if err != nil {
		return nil, err
	}

	return model.NewArticle(articleID, author, content), nil
}
