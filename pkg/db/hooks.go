package db

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

var _ bun.BeforeAppendModelHook = (*Fruit)(nil)

// BeforeAppendModel implements schema.BeforeAppendModelHook
func (f *Fruit) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		f.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		f.ModifiedAt = time.Now()
	}
	return nil
}
