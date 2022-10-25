package db

import (
	"fmt"
	"sort"
)

// fruitsCollection represent Fruits in DB
const fruitsCollection = "fruits"

// Fruit model to hold the Fruit data
type Fruit struct {
	ID     interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string      `json:"name" bson:"name" `
	Season string      `json:"season" bson:"season" `
	Emoji  string      `json:"emoji,omitempty" bson:"emoji,omitempty"`
}

// Fruits represents a collection of Fruits
type Fruits []*Fruit

var _ sort.Interface = (Fruits)(nil)

// Len implements sort.Interface
func (f Fruits) Len() int {
	return len(f)
}

// Less implements sort.Interface
func (f Fruits) Less(i int, j int) bool {
	return f[i].Name < f[j].Name
}

// Swap implements sort.Interface
func (f Fruits) Swap(i int, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Fruit) String() string {
	return fmt.Sprintf("ID: %v, Name: %s, Season: %s", f.ID, f.Name, f.Season)
}
