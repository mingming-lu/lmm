package simple

import (
	"context"
	"fmt"
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

	cli.Register("setCharacterSetTo-utf8mb4", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		rows, err := mysql.Query(c, "SHOW TABLES")
		if err != nil {
			return err
		}
		defer rows.Close()

		var table string

		for rows.Next() {
			if err := rows.Scan(&table); err != nil {
				return err
			}
			_, err := mysql.Exec(c, fmt.Sprintf("ALTER TABLE %s CONVERT TO CHARACTER SET utf8mb4;", table))
			if err != nil {
				return err
			}
		}

		return nil
	}))
}
