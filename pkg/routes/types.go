package routes

import (
	"github.com/kameshsampath/go-fruits-api/pkg/db"
)

//Endpoints is the marker interface for defining routes
type Endpoints struct {
	Config *db.Config
}

//NewEndpoints gives handle to REST Endpoints
//dbType could be one of "pg","mysql","sqlite".Defaults to "sqlite"
func NewEndpoints(dbc *db.Config) *Endpoints {
	return &Endpoints{
		Config: dbc,
	}
}
