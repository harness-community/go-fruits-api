package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/kameshsampath/go-fruits-api/pkg/db"
	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

//AddFruit godoc
// @Summary Add a fruit to Database
// @Description Adds a new Fruit to the Database
// @Tags fruit
// @Accept json
// @Produce json
// @Param message body db.Fruit true "Fruit object"
// @Success 200 {object} db.Fruit
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/add [post]
func (e *Endpoints) AddFruit(c echo.Context) error {
	log := e.Config.Log
	ctx := e.Config.Ctx
	dbConn := e.Config.DB
	f := &db.Fruit{}
	if err := c.Bind(f); err != nil {
		return err
	}
	log.Infof("Adding Fruit %s", f)
	err := dbConn.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := dbConn.NewInsert().
			Model(f).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Error adding fruit %v, %v", f, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Fruit %s successfully saved", f)
	c.JSON(http.StatusCreated, f)
	return nil
}

// DeleteFruit godoc
// @Summary Delete a fruit from Database
// @Description Deletes a Fruit to the Database
// @Tags fruit
// @Param id path int true "Fruit ID"
// @Success 204
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/{id} [delete]
func (e *Endpoints) DeleteFruit(c echo.Context) error {
	log := e.Config.Log
	ctx := e.Config.Ctx
	dbConn := e.Config.DB
	var ID int
	if err := echo.PathParamsBinder(c).
		Int("id", &ID).
		BindError(); err != nil {
		return err
	}
	if ID == 0 {
		err := fmt.Errorf("fruit with id %d not found", ID)
		utils.NewHTTPError(c, http.StatusNotFound, err)
		return err
	}
	f := &db.Fruit{
		ID: ID,
	}
	log.Infof("Deleting Fruit with id %d", ID)
	err := dbConn.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := dbConn.NewDelete().
			Model(f).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Error deleting fruit with ID %d, %v", ID, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Fruit with id  %d successfully deleted", ID)
	c.NoContent(http.StatusNoContent)
	return nil
}

// GetFruitsByName godoc
// @Summary Gets fruits by name
// @Description Gets list of fruits by name
// @Tags fruit
// @Produce json
// @Param name path string true "Full or partial name of the fruit"
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/{name} [get]
func (e *Endpoints) GetFruitsByName(c echo.Context) error {
	log := e.Config.Log
	var name string
	if err := echo.PathParamsBinder(c).
		String("name", &name).
		BindError(); err != nil {
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Getting Fruit with name %s", name)
	ctx := e.Config.Ctx
	dbConn := e.Config.DB
	var fruits = &db.Fruits{}
	if err := dbConn.NewSelect().
		Model(fruits).
		Where(`UPPER(name) LIKE ?`, fmt.Sprintf("%%%s%%", strings.ToUpper(name))).
		Scan(ctx); err != nil {
		log.Errorf("Error getting fruits by name %s, %v", name, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Found %d Fruits with name %s", fruits.Len(), name)
	c.JSON(http.StatusOK, fruits)
	return nil
}

// GetFruitsBySeason godoc
// @Summary Gets fruits by season
// @Description Gets a list of fruits by season
// @Tags fruit
// @Produce json
// @Param season path string true "Full or partial name of the season"
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/season/{season} [get]
func (e *Endpoints) GetFruitsBySeason(c echo.Context) error {
	log := e.Config.Log
	var season string
	if err := echo.PathParamsBinder(c).
		String("season", &season).
		BindError(); err != nil {
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Getting Fruit for season %s", season)
	ctx := e.Config.Ctx
	dbConn := e.Config.DB
	var fruits = &db.Fruits{}
	if err := dbConn.NewSelect().
		Model(fruits).
		Where("UPPER(season) = ?", strings.ToUpper(season)).
		Scan(ctx); err != nil {
		log.Errorf("Error getting fruits for season %s, %v", season, err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Found %d Fruits for season %s", fruits.Len(), season)
	c.JSON(http.StatusOK, fruits)
	return nil
}

// ListFruits godoc
// @Summary Gets all fruits
// @Description Gets a list all available fruits from the database
// @Tags fruit
// @Produce json
// @Success 200 {object} db.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits [get]
func (e *Endpoints) ListFruits(c echo.Context) error {
	log := e.Config.Log
	log.Infoln("Getting All Fruits ")
	ctx := e.Config.Ctx
	dbConn := e.Config.DB
	var fruits = &db.Fruits{}
	if err := dbConn.NewSelect().
		Model(fruits).
		Scan(ctx); err != nil {
		log.Errorf("Error getting all fruits, %v", err)
		utils.NewHTTPError(c, http.StatusInternalServerError, err)
		return err
	}
	log.Infof("Found %d Fruits", fruits.Len())
	c.JSON(http.StatusOK, fruits)
	return nil
}
