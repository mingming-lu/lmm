package persistence

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/service"
	"lmm/api/storage/db"
)

// ArticleStorage is a implementation of ArticleRepository
type ArticleStorage struct {
	db            db.DB
	authorService service.AuthorService
}

// NewArticleStorage constructs a new article repository with concrete struct
func NewArticleStorage(db db.DB, authorService service.AuthorService) *ArticleStorage {
	return &ArticleStorage{db: db, authorService: authorService}
}

// NextID generate a random string
func (s *ArticleStorage) NextID(c context.Context) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.New().String())))
}

// Save persist a article domain model
func (s *ArticleStorage) Save(c context.Context, article *model.Article) error {
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return err
	}

	findUserID, err := tx.Prepare(c, `
		select id from user where name = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	saveArticle, err := tx.Prepare(c, `
		insert into article (uid, alias_uid, user, title, body, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?)
		on duplicate key update id = LAST_INSERT_ID(id), alias_uid = ?, title = ?, body = ?, updated_at = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	deleteTags, err := tx.Prepare(c, `
		delete at from article_tag at left join article a on a.id = at.article where at.article = ?
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	saveTags, err := tx.Prepare(c, `
		insert into article_tag (article, sort, name) values (?, ?, ?)
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	now := time.Now()

	var userID int

	if err := findUserID.QueryRow(c, article.Author().Name()).Scan(&userID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	res, err := saveArticle.Exec(c,
		article.ID().Raw(),
		article.ID().String(),
		userID,
		article.Content().Text().Title(),
		article.Content().Text().Body(),
		now,
		now,
		article.ID().String(),
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

	if _, err := deleteTags.Exec(c, lastID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	for _, tag := range article.Content().Tags() {
		if _, err := saveTags.Exec(c, lastID, tag.ID().Order(), tag.Name()); err != nil {
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
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(c, `
		select id, uid, alias_uid, user, title, body from article where uid = ? for update
	`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	article, err := s.articleModelFromRow(c, stmt.QueryRow(c, id.String()))
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

func (s *ArticleStorage) articleModelFromRow(c context.Context, row *sql.Row) (*model.Article, error) {
	var (
		id             int
		rawArticleID   string
		aliasArticleID string
		userID         int
		title          string
		body           string
	)
	if err := row.Scan(&id, &rawArticleID, &aliasArticleID, &userID, &title, &body); err != nil {
		if err == db.ErrNoRows {
			return nil, domain.ErrNoSuchArticle
		}
		return nil, err
	}
	findUserName := s.db.Prepare(c, `select name from user where id = ?`)
	defer findUserName.Close()

	var userName string
	if err := findUserName.QueryRow(c, userID).Scan(&userName); err != nil {
		return nil, err
	}

	author, err := s.authorService.AuthorFromUserName(c, userName)
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

	stmt := s.db.Prepare(c, "select sort, name from article_tag where article = ?")
	defer stmt.Close()

	rows, err := stmt.Query(c, id)
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

	article := model.NewArticle(articleID, author, content)
	if err := article.ID().SetAlias(aliasArticleID); err != nil {
		return nil, err
	}
	return article, nil
}
