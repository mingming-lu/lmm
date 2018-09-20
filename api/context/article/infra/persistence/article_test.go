package persistence

import (
	"lmm/api/testing"
)

func TestSaveArticle_NewArticle(tt *testing.T) {
	t := testing.NewTester(tt)

	author, _ := authorService.AuthorFromUserID(user.ID())
	article, _ := articleService.NewArticleToPost(author, "title", "body", make([]string, 0))
	err := articleRepository.Save(article)
	t.NoError(err)
}
