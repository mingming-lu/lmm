package simple

import (
	"context"
	"lmm/api/storage/db"
	"log"

	"lmm/api/cli"
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
	cli.Register("addRoleColumnToUserTable", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		_, err := mysql.Exec(c, `alter table user add column role VARCHAR(31) NOT NULL`)
		return err
	}))
	cli.Register("createTableUserRoleChangeHistory", NewCommand(func(c context.Context) error {
		mysql := db.DefaultMySQL()
		defer mysql.Close()

		_, err := mysql.Exec(c, `
CREATE TABLE IF NOT EXISTS user_role_change_history (
	id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	operator BIGINT UNSIGNED NOT NULL, -- user.id
	operator_role VARCHAR(31) NOT NULL, -- user.role
	target_user BIGINT UNSIGNED NOT NULL, -- user.id
	target_user_role VARCHAR(31) NOT NULL, -- user.role
	target_role VARCHAR(31) NOT NULL, -- user.role
	changed_at DATETIME NOT NULL,
	PRIMARY KEY (id),
	INDEX changed_at (changed_at)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
`)
		return err
	}))
}
