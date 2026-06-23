package database

import (
	"context"

	"github.com/uptrace/bun"
)

func AutoMigration(db *bun.DB) {
	ctx := context.Background()

	models := []interface{}{}

	for _, m := range models {
		_, err := db.NewCreateTable().
			Model(m).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}
