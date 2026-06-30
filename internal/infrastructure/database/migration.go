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

	ensureAnomalyEventEnrichmentColumns(ctx, db)
}

func ensureAnomalyEventEnrichmentColumns(ctx context.Context, db *bun.DB) {
	statements := []string{
		`ALTER TABLE anomaly_events ADD COLUMN IF NOT EXISTS latency_sum_ms BIGINT NOT NULL DEFAULT 0 CHECK (latency_sum_ms >= 0)`,
		`ALTER TABLE anomaly_events ADD COLUMN IF NOT EXISTS avg_latency_ms DOUBLE PRECISION NOT NULL DEFAULT 0`,
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			panic(err)
		}
	}
}
