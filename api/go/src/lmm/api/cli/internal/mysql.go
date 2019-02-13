package internal

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"lmm/api/storage/db"

	"github.com/pkg/errors"
	"github.com/schemalex/schemalex"
	"github.com/schemalex/schemalex/diff"
)

const (
	mysqlSchemaFile = "/sql/schema.sql"
)

// MySQLSchemaDDL generates MySQL schema ddl, applies ddl if deploy is true
func MySQLSchemaDDL(c context.Context, deploy bool) error {
	schema, err := schemalex.NewSchemaSource("file://" + mysqlSchemaFile)
	if err != nil {
		return errors.Wrap(err, "cannot open /sql/schema.sql")
	}

	dsn := "mysql://" + db.DefaultMySQLConfig().DSN()
	database, err := schemalex.NewSchemaSource(dsn)
	if err != nil {
		return errors.Wrap(err, dsn)
	}

	dst := &bytes.Buffer{}
	parser := schemalex.New()

	if err := diff.Sources(dst, database, schema, diff.WithParser(parser)); err != nil {
		return err
	}

	dst.WriteByte('\n') // ensure ends with '\n' (for beautifully printing)

	if deploy {
		return deployMySQLSchemaDDL(c, strings.Split(dst.String(), ";"))
	}

	io.Copy(os.Stdout, dst)
	return nil
}

func deployMySQLSchemaDDL(c context.Context, stmts []string) error {
	mysql := db.DefaultMySQL()
	defer mysql.Close()

	tx, err := mysql.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		fmt.Println(stmt)
		if _, err := tx.Exec(c, stmt); err != nil {
			return db.RollbackWithError(tx, err)
		}
	}

	return tx.Commit()
}
