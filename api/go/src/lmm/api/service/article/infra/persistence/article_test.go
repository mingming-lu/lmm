package persistence

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

func TestNextArticleID(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	idStr := articleRepository.NextID(c)
	id, err := model.NewArticleID(idStr)
	if !t.NoError(err) {
		t.Logf(err.Error(), idStr)
	}
	t.NotNil(id)
	t.Is(idStr, id.String())
}

func TestSaveArticle(tt *testing.T) {
	c := context.Background()

	user := testutil.NewUser(mysql)

	author, err := authorService.AuthorFromUserName(c, user.Name())
	if err != nil {
		tt.Fatal(err)
	}

	tt.Run("Success", func(tt *testing.T) {
		cases := map[string]struct {
			Title string
			Body  string
			Tags  []string
		}{
			"SaveFirstArticle": {
				Title: uuid.New().String()[:8],
				Body:  uuid.New().String(),
				Tags:  randomStringSlice(),
			},
			"SaveSecondArticle": {
				Title: uuid.New().String()[:8],
				Body:  uuid.New().String(),
				Tags:  randomStringSlice(),
			},
		}

		for testName, testCase := range cases {
			tt.Run(testName, func(tt *testing.T) {
				t := testing.NewTester(tt)
				article, err := articleService.NewArticleToPost(c, author,
					testCase.Title,
					testCase.Body,
					testCase.Tags,
				)
				t.NoError(err)
				t.NoError(articleRepository.Save(c, article))

				id, title, body, err := selectArticleWhereUIDIs(article.ID().String())
				t.NoError(err)
				t.Is(testCase.Title, title)
				t.Is(testCase.Body, body)

				tags, err := selectTagsWhereArticleIDIs(id)
				t.NoError(err)
				t.Are(testCase.Tags, tags)

				tt.Run("EditArticle", func(tt *testing.T) {
					t := testing.NewTester(tt)
					text, err := model.NewText(uuid.New().String()[:8], uuid.New().String())
					t.NoError(err)
					tags := func() []*model.Tag {
						tags := make([]*model.Tag, 0)
						for i, tagName := range randomStringSlice() {
							tag, err := model.NewTag(article.ID(), uint(i+1), tagName)
							t.NoError(err)
							tags = append(tags, tag)
						}
						return tags
					}()
					content, err := model.NewContent(text, tags)
					t.NoError(err)

					article.EditContent(content)
					t.NoError(articleRepository.Save(c, article))

					id, title, body, err := selectArticleWhereUIDIs(article.ID().String())
					t.NoError(err)
					t.Is(text.Title(), title)
					t.Is(text.Body(), body)
					tagsGot, err := selectTagsWhereArticleIDIs(id)
					t.NoError(err)
					t.Are(func() []string {
						names := make([]string, 0)
						for _, tag := range tags {
							names = append(names, tag.Name())
						}
						return names
					}(), tagsGot)
				})
			})
		}
	})
}

func TestFindArticleByID(tt *testing.T) {
	c := context.Background()

	user := testutil.NewUser(mysql)

	author, err := authorService.AuthorFromUserName(c, user.Name())
	if err != nil {
		tt.Fatal(err)
	}

	article, err := articleService.NewArticleToPost(c, author, "awesome title", "awesome body", nil)
	if err != nil {
		tt.Fatal(err)
	}

	tt.Run("NotFound", func(tt *testing.T) {
		t := testing.NewTester(tt)
		articleGot, err := articleRepository.FindByID(c, article.ID())
		t.IsError(domain.ErrNoSuchArticle, err)
		t.Nil(articleGot)
	})

	if err := articleRepository.Save(c, article); err != nil {
		tt.Fatal(err)
	}

	tt.Run("Found", func(tt *testing.T) {
		t := testing.NewTester(tt)
		articleGot, err := articleRepository.FindByID(c, article.ID())
		t.NoError(err)
		t.Is(article, articleGot)
	})
}

func selectArticleWhereUIDIs(uid string) (int, string, string, error) {
	var (
		articleID int
		title     string
		body      string
	)

	err := mysql.QueryRow(context.Background(), `
		select id, title, body from article where uid = ?
	`, uid).Scan(&articleID, &title, &body)

	if err != nil {
		return 0, "", "", err
	}

	return articleID, title, body, nil
}

func selectTagsWhereArticleIDIs(id int) ([]string, error) {
	rows, err := mysql.Query(context.Background(), `
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

func randomStringSlice() []string {
	rand.Seed(time.Now().UnixNano())
	s := make([]string, 0)

	for i := 0; i < rand.Intn(10); i++ {
		s = append(s, uuid.New().String()[:8])
	}

	return s
}
