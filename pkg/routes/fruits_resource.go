package routes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const (
	DMLLISTFRUITS       = "SELECT * FROM fruits ORDER BY name ASC"
	DMLINSERTFRUIT      = "INSERT INTO fruits(name,season,emoji) values(?,?,?)"
	DMLGETFRUITBYNAME   = `SELECT * FROM fruits WHERE NAME LIKE ? COLLATE NOCASE ORDER BY name ASC`
	DMLGETFRUITBYSEASON = `SELECT * FROM fruits WHERE SEASON LIKE ? COLLATE NOCASE ORDER BY name ASC`
	DMLFRUITBYID        = "DELETE FROM fruits WHERE id = ?"
)

var (
	err error
)

type Fruit struct {
	Id     int64  `json:"id,omitempty" from:"id" uri:"id"`
	Name   string `json:"name" from:"name" uri:"name"`
	Season string `json:"season" from:"season" uri:"season"`
	Emoji  string `json:"emoji,omitempty" from:"emoji"`
}

type fruitsResponse []Fruit

//FruitsResource builds and handles GreetingResource URI a simple CRUD mapping to DB
// via /api/fruits for demonstration purpose
func FruitsResource(rg *gin.RouterGroup, db *sql.DB) {
	health := rg.Group("/api")

	health.POST("/fruits/add", addFruit(db))
	health.GET("/fruits", listFruits(db))
	health.DELETE("/fruits/:id", deleteFruit(db))
	health.GET("/fruits/:name", getFruitsByName(db))
	health.GET("/fruits/season/:season", getFruitsBySeason(db))
}

func addFruit(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fruit Fruit
		if err = c.ShouldBind(&fruit); err != nil {
			c.JSON(http.StatusNotFound, map[string]string{
				"reason": fmt.Sprintf("Error saving Fruit %s", err),
			})
		} else {
			log.Printf("Saving Fruit %v", fruit)
			if stmt, err := db.Prepare(DMLINSERTFRUIT); err != nil {
				c.JSON(http.StatusNotFound, map[string]string{
					"reason": fmt.Sprintf("Error saving Fruit %s", err),
				})
			} else {
				if fruit.Emoji == "" {
					//default set some plant emoji
					fruit.Emoji = "U+1F33F"
				}
				tx, err := db.Begin()
				if err != nil {
					c.JSON(http.StatusNotFound, map[string]string{
						"reason": fmt.Sprintf("Error saving Fruit %s", err),
					})
				} else {
					if rs, err := stmt.Exec(fruit.Name, fruit.Season, fruit.Emoji); err != nil {
						c.JSON(http.StatusNotFound, map[string]string{
							"reason": fmt.Sprintf("Error saving Fruit %s", err),
						})
						if err = tx.Rollback(); err != nil {
							log.Fatalf("Unable to rollback transaction %s", err)
						}
					} else {
						if pk, err := rs.LastInsertId(); err != nil {
							log.Println("Unable to get the primary key, rolling back transaction")
							c.JSON(http.StatusNotFound, nil)
							if err = tx.Rollback(); err != nil {
								log.Fatalf("Unable to rollback transaction %s", err)
							}
						} else {
							log.Printf("Successfully saved, with id %d", pk)
							c.JSON(http.StatusCreated, &map[string]int64{"id": pk})
						}
						tx.Commit()
					}
				}
			}
		}
	}
}

func deleteFruit(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if stmt, err := db.Prepare(DMLFRUITBYID); err != nil {
			c.JSON(http.StatusNotFound, map[string]string{
				"reason": fmt.Sprintf("Error deleting Fruit with id %d, %s", id, err),
			})
		} else {
			tx, err := db.Begin()
			if err != nil {
				c.JSON(http.StatusNotFound, map[string]string{
					"reason": fmt.Sprintf("Error deleting Fruit with id %d, %s", id, err),
				})
			} else {
				if _, err = stmt.Exec(id); err != nil {
					c.JSON(http.StatusNotFound, map[string]string{
						"reason": fmt.Sprintf("Error deleting Fruit with id %s, %s", id, err),
					})
					if err = tx.Rollback(); err != nil {
						log.Fatalf("Unable to rollback transaction %s", err)
					}
				} else {
					log.Printf("Successfully deleted fruit with id %s", id)
					c.Writer.WriteHeader(http.StatusNoContent)
					tx.Commit()
				}
			}
		}
	}
}

func getFruitsByName(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if stmt, err := db.Prepare(DMLGETFRUITBYNAME); err != nil {
			c.JSON(http.StatusNotFound, map[string]string{
				"reason": fmt.Sprintf("Error getting fruits by season %s, %s", name, err),
			})
		} else {
			if rows, err := stmt.Query("%" + name + "%"); err != nil {
				c.JSON(http.StatusNotFound, map[string]string{
					"reason": fmt.Sprintf("Error getting fruits by season %s, %s", name, err),
				})
			} else {
				fr := buildFruitsResponse(rows, err)
				c.JSON(http.StatusOK, fr)
			}
		}
	}
}

func getFruitsBySeason(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		season := c.Param("season")
		if stmt, err := db.Prepare(DMLGETFRUITBYSEASON); err != nil {
			c.JSON(http.StatusNotFound, map[string]string{
				"reason": fmt.Sprintf("Error getting fruits by season %s, %s", season, err),
			})
		} else {
			if rows, err := stmt.Query("%" + season + "%"); err != nil {
				c.JSON(http.StatusNotFound, map[string]string{
					"reason": fmt.Sprintf("Error getting fruits by season %s, %s", season, err),
				})
			} else {
				fr := buildFruitsResponse(rows, err)
				c.JSON(http.StatusOK, fr)
			}
		}
	}
}

func listFruits(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err == nil {
			if rows, err := db.Query(DMLLISTFRUITS); err != nil {
				c.JSON(http.StatusNotFound, map[string]string{
					"reason": fmt.Sprintf("Error querying fruits %s", err),
				})
			} else {
				fr := buildFruitsResponse(rows, err)
				c.JSON(http.StatusOK, fr)
			}
		} else {
			c.JSON(http.StatusNotFound, fmt.Sprintf("API unavailable %s", err))
		}
	}
}

func buildFruitsResponse(rows *sql.Rows, err error) fruitsResponse {
	defer rows.Close()
	var fr fruitsResponse
	for rows.Next() {
		var f Fruit
		if err = rows.Scan(&f.Id, &f.Name, &f.Season, &f.Emoji); err == nil {
			fr = append(fr, f)
		} else {
			log.Fatalf("Error reading row %s", err)
		}
	}
	return fr
}
