package routes

import (
	"fmt"
	"net/http"

	"github.com/kameshsampath/go-fruits-api/pkg/db"
	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddFruit godoc
// @Summary Add a fruit to Database
// @Description Adds a new Fruit to the Database
// @Tags fruit
// @Accept json
// @Produce json
// @Param message body db.Fruit true "Fruit object"
// @Success 200 {object} db.Fruit
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/add [post]
func (e *Endpoints) AddFruit(c echo.Context) error {
	log := e.Config.Log
	ctx := e.Config.Ctx
	database := e.Config.DB

	f := &db.Fruit{}
	if err := c.Bind(f); err != nil {
		return err
	}
	log.Infof("Adding Fruit %s", f)
	r, err := database.Collection(e.Config.Collection).InsertOne(ctx, f)
	if err != nil {
		log.Errorf("Error adding fruit %v, %v", f, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	f.ID = r.InsertedID
	log.Infof("Fruit %s successfully saved", f)
	return c.JSON(http.StatusCreated, f)
}

// DeleteAll godoc
// @Summary Deletes all the fruits from Database
// @Description Deletes all the fruits from Database
// @Tags fruit
// @Success 204
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/ [delete]
func (e *Endpoints) DeleteAll(c echo.Context) error {
	log := e.Config.Log
	ctx := e.Config.Ctx
	database := e.Config.DB
	_, err := database.Collection(e.Config.Collection).DeleteMany(ctx, bson.D{})
	if err != nil {
		log.Errorf("Error deleting all fruits,%v", err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infoln("All fruits deleted")
	return c.NoContent(http.StatusNoContent)
}

// DeleteFruit godoc
// @Summary Delete a fruit from Database
// @Description Deletes a Fruit to the Database
// @Tags fruit
// @Param id path string true "Fruit ID"
// @Success 204
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/{id} [delete]
func (e *Endpoints) DeleteFruit(c echo.Context) error {
	log := e.Config.Log
	ctx := e.Config.Ctx
	database := e.Config.DB
	var ID string
	if err := echo.PathParamsBinder(c).
		String("id", &ID).
		BindError(); err != nil {
		return err
	}
	if ID == "" {
		err := fmt.Errorf("fruit with id %s not found", ID)
		utils.NewHTTPError(c, http.StatusNotFound, err)
		return err
	}
	//make sure that we get only valid ObjectID for deletions
	//otherwise return the error
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Errorf("Error deleting fruit with ID %s, %v", ID, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	f := bson.D{{"_id", objID}} //nolint:govet
	log.Infof("Deleting Fruit with id %s", ID)
	_, err = database.Collection(e.Config.Collection).DeleteOne(ctx, f)
	if err != nil {
		log.Errorf("Error deleting fruit with ID %s, %v", ID, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Fruit with id %s successfully deleted", ID)
	return c.NoContent(http.StatusNoContent)
}

// GetFruitsByName godoc
// @Summary Gets fruits by name
// @Description Gets list of fruits by name
// @Tags fruit
// @Produce json
// @Param name path string true "Full or partial name of the fruit"
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/search/{name} [get]
func (e *Endpoints) GetFruitsByName(c echo.Context) error {
	return e.fruitFinder(c, "name")
}

// GetFruitsBySeason godoc
// @Summary Gets fruits by season
// @Description Gets a list of fruits by season
// @Tags fruit
// @Produce json
// @Param season path string true "Full or partial name of the season"
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/season/{season} [get]
func (e *Endpoints) GetFruitsBySeason(c echo.Context) error {
	return e.fruitFinder(c, "season")
}

// ListFruits godoc
// @Summary Gets all fruits
// @Description Gets a list all available fruits from the database
// @Tags fruit
// @Produce json
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
// @Router /fruits/ [get]
func (e *Endpoints) ListFruits(c echo.Context) error {
	log := e.Config.Log
	log.Infoln("Getting All Fruits ")
	ctx := e.Config.Ctx
	database := e.Config.DB
	cur, err := database.Collection(e.Config.Collection).
		Find(ctx, bson.D{},
			options.
				Find().
				SetSort(bson.D{{"name", 1}})) //nolint:govet
	if err != nil {
		log.Errorf("Error getting all fruits, %v", err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	var fruits = &db.Fruits{}
	err = cur.All(ctx, fruits)
	if err != nil {
		log.Errorf("Error getting all fruits, %v", err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Found %d Fruits", fruits.Len())
	return c.JSON(http.StatusOK, fruits)
}

func (e *Endpoints) fruitFinder(c echo.Context, filterBy string) error {
	log := e.Config.Log
	var filterV string
	if err := echo.PathParamsBinder(c).
		String(filterBy, &filterV).
		BindError(); err != nil {
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Getting Fruit with %s %s", filterBy, filterV)
	ctx := e.Config.Ctx
	database := e.Config.DB
	p := primitive.Regex{
		Pattern: filterV,
		Options: "i"}
	filter := bson.D{
		{ //nolint:govet
			filterBy,
			bson.D{
				{"$regex", p}, //nolint:govet
			},
		},
	}
	cur, err := database.Collection(e.Config.Collection).Find(ctx, filter)
	if err != nil {
		log.Errorf("Error getting fruits by %s %s, %v", filterBy, filterV, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}

	var fruits = &db.Fruits{}
	err = cur.All(ctx, fruits)
	if err != nil {
		log.Errorf("Error getting fruits by %s %s, %v", filterBy, filterV, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Found %d Fruits by %s %s", fruits.Len(), filterBy, filterV)
	return c.JSON(http.StatusOK, fruits)
}
