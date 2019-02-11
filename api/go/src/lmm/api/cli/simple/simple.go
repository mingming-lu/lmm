package simple

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"

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
	cli.Register("createUserPasswordChangeHistoryTable", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		_, err := mysql.Exec(c, `
CREATE TABLE IF NOT EXISTS user_password_change_history (
	id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	user BIGINT UNSIGNED NOT NULL, -- user.id
	changed_at DATETIME NOT NULL,
	PRIMARY KEY (id),
	INDEX user_change_history (user, changed_at)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
		`)
		return err
	}))

	cli.Register("addEmailColumnToUserTable", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		g := errgroup.Group{}
		g.Go(func() error {
			mysql.Exec(c, `ALTER TABLE user ADD COLUMN email VARCHAR(255) NOT NULL;`)
			return nil
		})

		g.Go(func() error {
			if _, err := mysql.Exec(c, `ALTER TABLE user drop KEY email`); err != nil {
				return err
			}

			if _, err := mysql.Exec(c, `ALTER TABLE user ADD UNIQUE KEY email (email);`); err != nil {
				return err
			}

			return nil
		})

		return g.Wait()
	}))
}
