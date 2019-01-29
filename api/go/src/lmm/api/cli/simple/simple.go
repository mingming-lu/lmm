package simple

import (
	"context"

	"lmm/api/cli"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/storage/db"
)

type NewCommand func(c context.Context) error

func (exec NewCommand) Execute(c context.Context) error {
	return exec(c)
}

func init() {
	cli.Register("resetAllArticlesUIDAndAliasUID", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		articleRepo := persistence.NewArticleStorage(mysql, nil)

		rows, err := mysql.Query(c, "select id from article")
		if err != nil {
			panic(err)
		}

		ids := make([]string, 0)

		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				panic(err)
			}
			ids = append(ids, id)
		}
		if err := rows.Close(); err != nil {
			panic(err)
		}

		stmt := mysql.Prepare(c, "update article set uid = ?, alias_uid = ? where id = ?")
		defer stmt.Close()

		for _, id := range ids {
			articleID := articleRepo.NextID(c)
			if _, err := stmt.Exec(c, articleID, articleID, id); err != nil {
				panic(err)
			}
		}

		return nil
	}))
}
