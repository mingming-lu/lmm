package simple

import (
	"context"
	"log"

	"lmm/api/cli"
	"lmm/api/cli/internal"
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

	cli.Register("MySQL-Schema-DDL-Dry-Run", NewCommand(func(c context.Context) error {
		return internal.MySQLSchemaDDL(c, false)
	}))

	cli.Register("MySQL-Schema-DDL-Deploy", NewCommand(func(c context.Context) error {
		return internal.MySQLSchemaDDL(c, true)
	}))
}
