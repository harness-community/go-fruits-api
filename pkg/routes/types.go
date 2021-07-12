package routes

import "database/sql"

//Endpoints is the marker interface for defining routes
type Endpoints struct {
	DB *sql.DB
}

//NewEndpoints gives handle to REST Endpoints
func NewEndpoints() *Endpoints {
	return &Endpoints{}
}
