package db

import (
	"context"
	"os"
	"testing"

	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func TestInitDB(t *testing.T) {
	log := utils.LogSetup(os.Stdout, "debug")
	dbc := New(
		WithContext(context.TODO()),
		WithDBFile("testdata/test.db"),
		WithLogger(log))

	dbc.Init()

	err := dbc.DB.Ping()

	if err != nil {
		t.Fatal(err)
	}

	dbfx := dbfixture.New(dbc.DB, dbfixture.WithRecreateTables())
	if err := dbfx.Load(dbc.Ctx, os.DirFS("."), "testdata/fixtures.yaml"); err != nil {
		t.Fatal(err)
	}

	expected := dbfx.MustRow("Fruit.mango").(*Fruit)
	assert.NotNil(t, expected)

	actual := &Fruit{}

	err = dbc.DB.NewSelect().
		Model(actual).
		Where("? = ?", bun.Ident("name"), "Mango").
		Scan(dbc.Ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, 1, "Expected ID to be  %d but got %d", 1, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "Expected Name to be  %s but got %s", expected.Name, actual.Name)
	assert.Equal(t, expected.Name, actual.Name, "Expected Season to be  %s but got %s", expected.Season, actual.Season)

	tearDown()
}

func tearDown() {
	os.Remove("testdata/test.db")
}
