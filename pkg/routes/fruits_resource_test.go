package routes

import (
	"context"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kameshsampath/go-fruits-api/pkg/db"
	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
)

var (
	testData = []interface{}{
		bson.D{{"_id", "1"}, {"name", "Mango"}, {"season", "Spring"}, {"emoji", "U+1F96D"}},      //nolint:govet
		bson.D{{"_id", "2"}, {"name", "Strawberry"}, {"season", "Spring"}, {"emoji", "U+1F96D"}}, //nolint:govet
		bson.D{{"_id", "3"}, {"name", "Orange"}, {"season", "Winter"}, {"emoji", "U+1F34B"}},     //nolint:govet
		bson.D{{"_id", "4"}, {"name", "Lemon"}, {"season", "Winter"}, {"emoji", "U+1F34A"}},      //nolint:govet
		bson.D{{"_id", "5"}, {"name", "Blueberry"}, {"season", "Summer"}, {"emoji", "U+1FAD0"}},  //nolint:govet
		bson.D{{"_id", "6"}, {"name", "Banana"}, {"season", "Summer"}, {"emoji", "U+1F34C"}},     //nolint:govet
		bson.D{{"_id", "7"}, {"name", "Watermelon"}, {"season", "Summer"}, {"emoji", "U+1F349"}}, //nolint:govet
		bson.D{{"_id", "8"}, {"name", "Apple"}, {"season", "Fall"}, {"emoji", "U+1F34E"}},        //nolint:govet
		bson.D{{"_id", "9"}, {"name", "Pear"}, {"season", "Fall"}, {"emoji", "U+1F350"}},         //nolint:govet
	}
	log *logrus.Logger
)

func loadFixtures() (*db.Config, error) {
	log = utils.LogSetup(os.Stdout, utils.LookupEnvOrString("TEST_LOG_LEVEL", "info"))
	dbc := db.New(utils.LookupEnvOrString("QUARKUS_MONGODB_CONNECTION_STRING", "mongodb://localhost:27017"),
		db.WithContext(context.TODO()),
		db.WithLogger(log),
		db.WithCollection(utils.LookupEnvOrString("FRUITS_DB_COLLECTION", "fruits")),
		db.WithDB(utils.LookupEnvOrString("FRUITS_DB", "testdb")))
	err := dbc.Init()
	if err != nil {
		return nil, err
	}

	//Clear existing data before starting any new tests
	dbc.DB.Collection(dbc.Collection).Drop(dbc.Ctx)
	_, err = dbc.DB.Collection(dbc.Collection).InsertMany(dbc.Ctx, testData)

	if err != nil {
		return nil, err
	}

	return dbc, nil
}

func TestAddFruit(t *testing.T) {
	dbc, err := loadFixtures()
	if err != nil {
		t.Fatal(err)
	}
	testCases := map[string]struct {
		requestBody string
		statusCode  int
		want        db.Fruit
	}{
		"withId": {
			requestBody: `{
        "_id": "10",
        "name": "Test Fruit",
        "season": "Summer"
        }`,
			statusCode: http.StatusCreated,
			want: db.Fruit{
				ID:     "10",
				Name:   "Test Fruit",
				Season: "Summer",
			},
		},
		"withoutId": {
			requestBody: `{
        "name": "Test Fruit 2",
        "season": "Spring"
        }`,
			statusCode: http.StatusCreated,
			want: db.Fruit{
				Name:   "Test Fruit 2",
				Season: "Spring",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var got db.Fruit
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/fruits/add", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ep := &Endpoints{
				Config: dbc,
			}
			if c := e.NewContext(req, rec); assert.NoError(t, ep.AddFruit(c)) {
				assert.Equal(t, tc.statusCode, rec.Code)
				database := ep.Config.DB
				ctx := context.TODO()
				opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(1) //nolint:govet
				cur, err := database.Collection(ep.Config.Collection).Find(ctx, bson.D{}, opts)
				if err != nil {
					t.Fatal(err)
				}

				if tc.want.ID == nil {
					var results []bson.D
					err = cur.All(ctx, &results)
					if err != nil {
						t.Fatal(err)
					}
					tc.want.ID = results[0].Map()["_id"]
				}
				res := database.Collection(ep.Config.Collection).FindOne(ctx, bson.D{{"name", tc.want.Name}}) //nolint:govet
				err = res.Decode(&got)
				if err != nil {
					t.Fatal(err)
				}
				//Verify Fruit
				if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(db.Fruit{})); diff != "" {
					t.Errorf("AddFruit() mismatch (-want +got):\n%s", diff)
				}
				//delete the added fruit
				database.Collection(ep.Config.Collection).DeleteOne(ctx, bson.D{{"name", got.Name}}) //nolint:govet
			}
		})
	}
}

func TestDeleteFruit(t *testing.T) {
	fruitID := "5"
	dbc, err := loadFixtures()
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/fruits/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/fruits/:id")
	c.SetParamNames("id")
	c.SetParamValues(fruitID)
	ep := &Endpoints{
		Config: dbc,
	}
	if assert.NoError(t, ep.DeleteFruit(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		database := ep.Config.DB
		ctx := context.TODO()
		got := &db.Fruit{}
		err = database.
			Collection("fruits").
			FindOne(ctx, bson.D{{"_id", fruitID}}). //nolint:govet
			Decode(got)
		assert.EqualError(t, err, "mongo: no documents in result", "Expecting Fruit with ID %s not to exist but it does", fruitID)
	}
}

func TestGetFruitByName(t *testing.T) {
	testCases := map[string]struct {
		name string
		want db.Fruits
	}{
		"default": {
			name: "Apple",
			want: db.Fruits{
				{
					ID:     "8",
					Name:   "Apple",
					Emoji:  "U+1F34E",
					Season: "Fall",
				},
			},
		},
		"allLower": {
			name: "apple",
			want: db.Fruits{
				{
					ID:     "8",
					Name:   "Apple",
					Emoji:  "U+1F34E",
					Season: "Fall",
				},
			},
		},
		"allCaps": {
			name: "APPLE",
			want: db.Fruits{
				{
					ID:     "8",
					Name:   "Apple",
					Emoji:  "U+1F34E",
					Season: "Fall",
				},
			},
		},
		"mixedCase": {
			name: "apPLe",
			want: db.Fruits{
				{
					ID:     "8",
					Name:   "Apple",
					Emoji:  "U+1F34E",
					Season: "Fall",
				},
			},
		},
	}
	dbc, err := loadFixtures()
	if err != nil {
		t.Fatal(err)
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/fruits/search/:name", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/fruits/search/:name")
			c.SetParamNames("name")
			c.SetParamValues(tc.name)
			ep := &Endpoints{
				Config: dbc,
			}
			if assert.NoError(t, ep.GetFruitsByName(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				var got db.Fruits
				b := rec.Body.Bytes()
				err := json.Unmarshal(b, &got)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, got, "Expecting the response to have Fruit(s) object but got none")
				//Verify Fruit
				if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(db.Fruit{})); diff != "" {
					t.Errorf("GetFruitsByName() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetFruitsBySeason(t *testing.T) {
	testCases := map[string]struct {
		season string
		want   db.Fruits
	}{
		"default": {
			season: "Summer",
			want: db.Fruits{
				{
					ID:     "5",
					Name:   "Blueberry",
					Emoji:  "U+1FAD0",
					Season: "Summer",
				},
				{
					ID:     "6",
					Name:   "Banana",
					Emoji:  "U+1F34C",
					Season: "Summer",
				},
				{
					ID:     "7",
					Name:   "Watermelon",
					Emoji:  "U+1F349",
					Season: "Summer",
				},
			},
		},
		"allLower": {
			season: "summer",
			want: db.Fruits{
				{
					ID:     "5",
					Name:   "Blueberry",
					Emoji:  "U+1FAD0",
					Season: "Summer",
				},
				{
					ID:     "6",
					Name:   "Banana",
					Emoji:  "U+1F34C",
					Season: "Summer",
				},
				{
					ID:     "7",
					Name:   "Watermelon",
					Emoji:  "U+1F349",
					Season: "Summer",
				},
			},
		},
		"allCaps": {
			season: "SUMMER",
			want: db.Fruits{
				{
					ID:     "5",
					Name:   "Blueberry",
					Emoji:  "U+1FAD0",
					Season: "Summer",
				},
				{
					ID:     "6",
					Name:   "Banana",
					Emoji:  "U+1F34C",
					Season: "Summer",
				},
				{
					ID:     "7",
					Name:   "Watermelon",
					Emoji:  "U+1F349",
					Season: "Summer",
				},
			},
		},
		"mixedCase": {
			season: "suMMEr",
			want: db.Fruits{
				{
					ID:     "5",
					Name:   "Blueberry",
					Emoji:  "U+1FAD0",
					Season: "Summer",
				},
				{
					ID:     "6",
					Name:   "Banana",
					Emoji:  "U+1F34C",
					Season: "Summer",
				},
				{
					ID:     "7",
					Name:   "Watermelon",
					Emoji:  "U+1F349",
					Season: "Summer",
				},
			},
		},
	}
	dbc, err := loadFixtures()
	if err != nil {
		t.Fatal(err)
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/fruits/season/:season", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/fruits/season/:season")
			c.SetParamNames("season")
			c.SetParamValues(tc.season)
			ep := &Endpoints{
				Config: dbc,
			}
			if assert.NoError(t, ep.GetFruitsBySeason(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				var got db.Fruits
				b := rec.Body.Bytes()
				err := json.Unmarshal(b, &got)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, got, "Expecting the response to have Fruits object but got none")
				//Verify Fruit
				if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(db.Fruit{})); diff != "" {
					t.Errorf("GetFruitsBySeason() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetAllFruits(t *testing.T) {
	dbc, err := loadFixtures()
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/fruits", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ep := &Endpoints{
		Config: dbc,
	}
	var want db.Fruits
	cwd, _ := os.Getwd()
	bw, err := os.ReadFile(path.Join(cwd, "testdata", "list.json"))
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bw, &want)
	if err != nil {
		t.Fatal(err)
	}
	sort.Sort(want)
	if c := e.NewContext(req, rec); assert.NoError(t, ep.ListFruits(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var got db.Fruits
		b := rec.Body.Bytes()
		err := json.Unmarshal(b, &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, got, "Expecting the response to have Fruits object but got none")
		sort.Sort(got)
		//Verify Fruits
		if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(db.Fruit{})); diff != "" {
			t.Errorf("GetAllFruits() mismatch (-want +got):\n%s", diff)
		}
	}
}
