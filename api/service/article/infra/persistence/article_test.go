package persistence

import (
	"context"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/testing"
)

func TestSaveArticle(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	author, err := authorService.AuthorFromUserID(c, user.ID())
	t.NoError(err)

	tagNames := make([]string, 3)
	tagNames[0], tagNames[1], tagNames[2] = "1", "2", "3"
	article, err := articleService.NewArticleToPost(c, author, "title", "body", tagNames)
	t.NoError(err)

	article2, err := articleService.NewArticleToPost(c, author, "title", "body", tagNames)
	t.NoError(err)

	// save new article
	t.NoError(articleRepository.Save(c, article))
	t.NoError(articleRepository.Save(c, article2))

	articleID, title, body, err := selectArticleWhereUIDIs(article.ID().String())
	t.NoError(err)

	t.Is("title", title)
	t.Is("body", body)

	tagNamesGot, err := selectTagsWhereArticleIDIs(articleID)
	t.NoError(err)
	t.Are(tagNames, tagNamesGot)

	// edit content
	text, err := model.NewText("new title", "new body")
	t.NoError(err)
	tags := make([]*model.Tag, 0)
	tag1, err := model.NewTag(article.ID(), 1, "111")
	t.NoError(err)
	tag2, err := model.NewTag(article.ID(), 2, "222")
	t.NoError(err)
	tag3, err := model.NewTag(article.ID(), 3, "333")
	t.NoError(err)
	tags = append(tags, tag1, tag2, tag3)
	content, err := model.NewContent(text, tags)
	t.NoError(err)
	article.EditContent(content)

	// save updated article
	t.NoError(articleRepository.Save(c, article))

	_, title, body, err = selectArticleWhereUIDIs(article.ID().String())
	t.NoError(err)
	t.Is("new title", title)
	t.Is("new body", body)

	tagNamesGot, err = selectTagsWhereArticleIDIs(articleID)
	t.NoError(err)
	t.Is("111", tagNamesGot[0])
	t.Is("222", tagNamesGot[1])
	t.Is("333", tagNamesGot[2])
}

func TestFindArticleByID(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	author, err := authorService.AuthorFromUserID(c, user.ID())
	t.NoError(err)

	tagNames := make([]string, 3)
	tagNames[0], tagNames[1], tagNames[2] = "awesome2", "awesome3", "awesome1"
	article, err := articleService.NewArticleToPost(c, author, "awesome title", "awesome body", tagNames)
	t.NoError(err)

	t.NoError(articleRepository.Save(c, article))

	articleGot, err := articleRepository.FindByID(c, article.ID())
	t.NoError(err)
	t.Is(article, articleGot)
}

func TestFindArticleByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := model.NewArticleID("notexist")
	t.NoError(err)

	article, err := articleRepository.FindByID(context.Background(), articleID)
	t.IsError(domain.ErrNoSuchArticle, err)
	t.Nil(article)
}

func selectArticleWhereUIDIs(uid string) (int, string, string, error) {
	var (
		articleID int
		title     string
		body      string
	)

	err := testing.DB().QueryRow(`
		select id, title, body from article where uid = ?
	`, uid).Scan(&articleID, &title, &body)

	if err != nil {
		return 0, "", "", err
	}

	return articleID, title, body, nil
}

func selectTagsWhereArticleIDIs(id int) ([]string, error) {
	rows, err := testing.DB().Query(`
		select name from article_tag where article = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tagName string
	tagNames := make([]string, 0)
	for rows.Next() {
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}
		tagNames = append(tagNames, tagName)
	}

	return tagNames, nil
}
