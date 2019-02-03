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
	cli.Register("addUserDescriptionIndex", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()

		_, err := mysql.Exec(c, `ALTER TABLE user ADD INDEX description (name, role, created_at)`)
		return err
	}))
}
