package data

const (
	//DDLFRUITSTABLE  creates the database
	DDLFRUITSTABLE = `
DROP TABLE IF EXISTS fruits;
CREATE TABLE IF NOT EXISTS fruits (
id SERIAL PRIMARY KEY ,
name VARCHAR NOT NULL,
season VARCHAR NOT NULL,
emoji VARCHAR)`
	DMLLISTFRUITS       = `SELECT * FROM fruits ORDER BY name ASC`
	DMLINSERTFRUIT      = `INSERT INTO fruits(name,season,emoji) values($1,$2,$3)`
	DMLGETFRUITBYNAME   = `SELECT * FROM fruits WHERE NAME LIKE ? COLLATE NOCASE ORDER BY name ASC`
	DMLGETFRUITBYSEASON = `SELECT * FROM fruits WHERE SEASON LIKE ? COLLATE NOCASE ORDER BY name ASC`
	DMLFRUITBYID        = `DELETE FROM fruits WHERE id = $1`
)

//Fruit model to hold the Fruit data
type Fruit struct {
	Id     int64  `json:"id,omitempty" from:"id" uri:"id"`
	Name   string `json:"name" from:"name" uri:"name"`
	Season string `json:"season" from:"season" uri:"season"`
	Emoji  string `json:"emoji,omitempty" from:"emoji"`
}

//Fruits represents a collection of Fruits
type Fruits []Fruit
