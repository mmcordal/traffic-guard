package database

import (
	"context"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

func AutoMigration(db *bun.DB) {
	ctx := context.Background()

	models := []interface{}{
		(*model.TrafficLog)(nil),
		(*model.TrafficBucket)(nil),
		(*model.DomainAnomalyCheck)(nil),
		(*model.AnomalyEvent)(nil),
	}

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
