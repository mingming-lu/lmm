package simple

import (
	"context"
	"log"

	"lmm/api/cli"
	"lmm/api/storage/db"
)

// NewCommand creates a cli.Command implement
type NewCommand func(c context.Context) error

// Execute implements cli.Command.Execute
func (exec NewCommand) Execute(c context.Context) error {
	return exec(c)
}

func init() {
	cli.Register("hello-world", NewCommand(func(_ context.Context) error {
		log.Println("hello world")
		return nil
	}))
	cli.Register("alterArticleAliasUIDUnique", NewCommand(func(c context.Context) error {
		db := db.DefaultMySQL()
		defer db.Close()

		_, err := db.Exec(c, "alter table article add unique key `alias_uid` (`alias_uid`)")
		return err
	}))
}
