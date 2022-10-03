package db

import (
	"fmt"
	"sort"
	"time"

	"github.com/uptrace/bun"
)

//Fruit model to hold the Fruit data
type Fruit struct {
	bun.BaseModel `bun:"table:fruits,alias:f"`

	ID         int       `bun:",pk,autoincrement" json:"id"`
	Name       string    `bun:",notnull" json:"name" `
	Season     string    `bun:",notnull" json:"season"`
	Emoji      string    `bun:"," json:"emoji,omitempty"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"-"`
	ModifiedAt time.Time `json:"-"`
}

//Fruits represents a collection of Fruits
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
	return fmt.Sprintf("ID: %d, Name: %s, Season: %s", f.ID, f.Name, f.Season)
}
