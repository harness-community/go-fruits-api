package db

import (
	"context"
	"os"
	"testing"

	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect"
)

func TestInitDB(t *testing.T) {
	log := utils.LogSetup(os.Stdout, utils.LookupEnvOrString("TEST_LOG_LEVEL", "info"))
	ctx := context.TODO()
	dbc := New(
		WithContext(ctx),
		WithLogger(log),
		WithDBType(utils.LookupEnvOrString("FRUITS_DB_TYPE", "sqlite")),
		WithDBFile("testdata/test.db"))
	dbc.Init()

	err := dbc.DB.Ping()

	if err != nil {
		t.Fatal(err)
	}

	dbc.DB.RegisterModel((*Fruit)(nil))

	dbfx := dbfixture.New(dbc.DB, dbfixture.WithRecreateTables())
	if err := dbfx.Load(dbc.Ctx, os.DirFS("."), "testdata/fixtures.yaml"); err != nil {
		t.Fatalf("Unable to load fixtures, %s", err)
	}

	expected := dbfx.MustRow("Fruit.mango").(*Fruit)
	assert.NotNil(t, expected)

	actual := &Fruit{}

	err = dbc.DB.NewSelect().
		Model(actual).
		Where("? = ?", bun.Ident("name"), "Mango").
		Scan(dbc.Ctx)
	assert.NoError(t, err)

	assert.Equal(t, 1, 1, "Expected ID to be  %d but got %d", 1, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "Expected Name to be  %s but got %s", expected.Name, actual.Name)
	assert.Equal(t, expected.Name, actual.Name, "Expected Season to be  %s but got %s", expected.Season, actual.Season)

	var lastID int
	var seqQuery string
	switch dbc.DBType {
	case dialect.PG:
		seqQuery = "SELECT currval(pg_get_serial_sequence('fruits','id'))"
	case dialect.MySQL:
		seqQuery = "SELECT LAST_INSERT_ID()"
	default:
		seqQuery = "SELECT ROWID from FRUITS order by ROWID DESC limit 1"
	}

	err = dbc.DB.NewRaw(seqQuery).Scan(dbc.Ctx, &lastID)
	assert.NoError(t, err)

	assert.Equal(t, 1, 1, "Expected Last Sequential ID to be  %d but got %d", 9, lastID)

	tearDown()
}

func tearDown() {
	os.Remove("testdata/test.db")
}
