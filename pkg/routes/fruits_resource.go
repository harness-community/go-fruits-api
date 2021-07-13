package routes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kameshsampath/gloo-fruits-api/pkg/data"
	"github.com/kameshsampath/gloo-fruits-api/pkg/utils"
	"log"
	"net/http"
	"strings"
)

var (
	err error
)

//AddFruit godoc
// @Summary Add fruit to Database
// @Description Adds a new Fruit to the Database
// @Tags fruit
// @Accept json
// @Produce json
// @Param message body data.Fruit true "Fruit object"
// @Success 200 {object} data.Fruit
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/add [post]
func (e *Endpoints) AddFruit(c *gin.Context) {
	var fruit data.Fruit
	if err = c.ShouldBind(&fruit); err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		return
	} else {
		log.Printf("Saving Fruit %v", fruit)
		if stmt, err := e.DB.Prepare(data.DMLINSERTFRUIT); err != nil {
			utils.NewError(c, http.StatusNotFound, err)
			return
		} else {
			if fruit.Emoji == "" {
				//default set some plant emoji
				fruit.Emoji = "U+1F33F"
			}
			if tx, err := e.DB.Begin(); err != nil {
				utils.NewError(c, http.StatusNotFound, err)
				return
			} else {
				if _, err := stmt.Exec(fruit.Name, fruit.Season, fruit.Emoji); err != nil {
					utils.NewError(c, http.StatusNotFound, err)
					if err = tx.Rollback(); err != nil {
						log.Fatalf("Unable to rollback transaction %s", err)
					}
					return
				} else {
					if rows, err := e.DB.Query(data.FRUITSIDSEQ); err != nil {
						log.Printf("Unable to get the primary key %s, rolling back transaction", err)
						utils.NewError(c, http.StatusNotFound, err)
						if err = tx.Rollback(); err != nil {
							log.Fatalf("Unable to rollback transaction %s", err)
						}
						return
					} else {
						for rows.Next() {
							if err = rows.Scan(&fruit.Id); err != nil {
								utils.NewError(c, http.StatusNotFound, err)
								if err = tx.Rollback(); err != nil {
									log.Fatalf("Unable to rollback transaction %s", err)
								}
								return
							}
							log.Printf("Successfully saved, with id %d", fruit.Id)
							c.JSON(http.StatusCreated, fruit)
							if err := tx.Commit(); err != nil {
								log.Fatalf("Unable to commit transaction %s", err)
							}
						}
					}
				}
			}
		}
	}
}

// DeleteFruit godoc
// @Summary Delete a fruit from Database
// @Description Deletes a Fruit to the Database
// @Tags fruit
// @Param fruit_id param path int true "Fruit ID"
// @Success 204
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/{fruit_id} [delete]
func (e *Endpoints) DeleteFruit(c *gin.Context) {
	id := c.Param("id")
	if stmt, err := e.DB.Prepare(data.DMLFRUITBYID); err != nil {
		log.Printf("Error deleting row %v, %s", stmt, err)
		utils.NewError(c, http.StatusNotFound, err)
		return
	} else {
		if tx, err := e.DB.Begin(); err != nil {
			utils.NewError(c, http.StatusNotFound, err)
			return
		} else {
			if _, err = stmt.Exec(id); err != nil {
				log.Printf("Error deleting row %v, %s", stmt, err)
				utils.NewError(c, http.StatusNotFound, err)
				if err = tx.Rollback(); err != nil {
					log.Fatalf("Unable to rollback transaction %s", err)
				}
				return
			} else {
				log.Printf("Successfully deleted fruit with id %s", id)
				c.Writer.WriteHeader(http.StatusNoContent)
				if err := tx.Commit(); err != nil {
					log.Fatalf("Unable to commit transaction %s", err)
				}
			}
		}
	}
}

// GetFruitsByName godoc
// @Summary Gets fruits by name
// @Description Gets list of fruits by name
// @Tags fruit
// @Produce json
// @Param fruit_name param path string true "Full or partial name of the fruit"
// @Success 200 {object} data.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/{fruit_name} [get]
func (e *Endpoints) GetFruitsByName(c *gin.Context) {
	name := c.Param("name")
	if stmt, err := e.DB.Prepare(data.DMLGETFRUITBYNAME); err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		return
	} else {
		if rows, err := stmt.Query("%" + strings.ToLower(name) + "%"); err != nil {
			utils.NewError(c, http.StatusNotFound, err)
			return
		} else {
			fr := buildFruitsResponse(rows, err)
			c.JSON(http.StatusOK, fr)
		}
	}
}

// GetFruitsBySeason godoc
// @Summary Gets fruits by season
// @Description Gets a list of fruits by season
// @Tags fruit
// @Produce json
// @Param season param path string true "Full or partial name of the season"
// @Success 200 {object} data.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits/season/{season} [get]
func (e *Endpoints) GetFruitsBySeason(c *gin.Context) {
	season := c.Param("season")
	if stmt, err := e.DB.Prepare(data.DMLGETFRUITBYSEASON); err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		return
	} else {
		if rows, err := stmt.Query("%" + strings.ToLower(season) + "%"); err != nil {
			utils.NewError(c, http.StatusNotFound, err)
			return
		} else {
			fr := buildFruitsResponse(rows, err)
			log.Printf("ROWS:%s", fr)
			c.JSON(http.StatusOK, fr)
		}
	}
}

// ListFruits godoc
// @Summary Gets all fruits
// @Description Gets a list all available fruits from the database
// @Tags fruit
// @Produce json
// @Success 200 {object} data.Fruits
// @Failure 404 {object} utils.HTTPError
//@Router /fruits [get]
func (e *Endpoints) ListFruits(c *gin.Context) {
	if err == nil {
		if rows, err := e.DB.Query(data.DMLLISTFRUITS); err != nil {
			utils.NewError(c, http.StatusNotFound, err)
			return
		} else {
			fr := buildFruitsResponse(rows, err)
			c.JSON(http.StatusOK, fr)
		}
	} else {
		c.JSON(http.StatusNotFound, fmt.Sprintf("API unavailable %s", err))
	}
}

func buildFruitsResponse(rows *sql.Rows, err error) data.Fruits {
	defer rows.Close()
	var fr data.Fruits
	for rows.Next() {
		var f data.Fruit
		if err = rows.Scan(&f.Id, &f.Name, &f.Season, &f.Emoji); err == nil {
			fr = append(fr, f)
		} else {
			log.Fatalf("Error reading row %s", err)
		}
	}
	return fr
}
