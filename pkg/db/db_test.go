package db

import (
	"context"
	"os"
	"testing"

	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	log := utils.LogSetup(os.Stdout, utils.LookupEnvOrString("TEST_LOG_LEVEL", "info"))
	dbc := New(utils.LookupEnvOrString("QUARKUS_MONGODB_CONNECTION_STRING", "mongodb://localhost:27017"),
		WithContext(context.TODO()),
		WithLogger(log),
		WithCollection(utils.LookupEnvOrString("FRUITS_DB_COLLECTION", "fruits")),
		WithDB(utils.LookupEnvOrString("FRUITS_DB", "testdb")))
	err := dbc.Init()
	require.NoError(t, err)
	go func() {
		defer dbc.Client.Disconnect(dbc.Ctx)
	}()
}
