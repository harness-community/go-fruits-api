package db

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

var _ bun.BeforeAppendModelHook = (*Fruit)(nil)

// BeforeAppendModel implements schema.BeforeAppendModelHook
func (m *Fruit) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.ModifiedAt = time.Now()
	}
	return nil
}
