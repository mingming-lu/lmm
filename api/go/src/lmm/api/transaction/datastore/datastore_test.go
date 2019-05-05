package datastore

import (
	"context"
	"lmm/api/testing"
	"os"

	"cloud.google.com/go/datastore"
)

func TestTxManager(tt *testing.T) {
	t := testing.NewTester(tt)

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if !t.NoError(err) {
		t.Fatalf(`failed to setup datastore: "%s"`, err.Error())
	}
	txManager := NewTransactionManager(client)

	tt.Run("New", func(tt *testing.T) {
		txCtx, err := txManager.New(ctx)
		t.NoError(err)
		t.NotNil(txCtx)

		tx, err := FromContext(txCtx)
		t.NoError(err)
		t.NotNil(tx)
	})
}
