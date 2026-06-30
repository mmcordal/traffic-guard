package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type AnomalyRepository interface {
	CreateAnomalyEvent(ctx context.Context, event *model.AnomalyEvent) error
	GetSourceIPBucketEnrichment(ctx context.Context, domain, sourceIP string, start, end time.Time) (*model.AnomalyEventEnrichment, error)
	GetEventsWithStartAndEndTime(ctx context.Context, domain string, start, end time.Time) ([]*model.AnomalyEvent, error)
}

type anomalyRepository struct {
	db *bun.DB
}

func NewAnomalyRepository(db *bun.DB) AnomalyRepository {
	return &anomalyRepository{db: db}
}

func (r *anomalyRepository) CreateAnomalyEvent(ctx context.Context, event *model.AnomalyEvent) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	return err
}

func (r *anomalyRepository) GetSourceIPBucketEnrichment(ctx context.Context, domain, sourceIP string, start, end time.Time) (*model.AnomalyEventEnrichment, error) {
	enrichment := &model.AnomalyEventEnrichment{}

	var dominant struct {
		QueryType model.QueryType `bun:"query_type"`
	}
	err := r.db.NewSelect().
		Model((*model.TrafficLog)(nil)).
		Column("query_type").
		Where("domain = ?", domain).
		Where("source_ip = ?", sourceIP).
		Where("timestamp >= ?", start).
		Where("timestamp < ?", end).
		Group("query_type").
		OrderExpr("COUNT(*) DESC").
		Order("query_type ASC").
		Limit(1).
		Scan(ctx, &dominant)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		enrichment.QueryType = dominant.QueryType
	}

	var latency struct {
		LatencySumMs sql.NullInt64   `bun:"latency_sum_ms"`
		AvgLatencyMs sql.NullFloat64 `bun:"avg_latency_ms"`
	}
	err = r.db.NewSelect().
		Model((*model.TrafficLog)(nil)).
		ColumnExpr("COALESCE(SUM(latency_ms), 0) AS latency_sum_ms").
		ColumnExpr("COALESCE(AVG(latency_ms), 0) AS avg_latency_ms").
		Where("domain = ?", domain).
		Where("source_ip = ?", sourceIP).
		Where("timestamp >= ?", start).
		Where("timestamp < ?", end).
		Scan(ctx, &latency)
	if err != nil {
		return nil, err
	}
	if latency.LatencySumMs.Valid {
		enrichment.LatencySumMs = latency.LatencySumMs.Int64
	}
	if latency.AvgLatencyMs.Valid {
		enrichment.AvgLatencyMs = latency.AvgLatencyMs.Float64
	}

	return enrichment, nil
}

func (r *anomalyRepository) GetEventsWithStartAndEndTime(ctx context.Context, domain string, start, end time.Time) ([]*model.AnomalyEvent, error) {
	var events []*model.AnomalyEvent
	err := r.db.NewSelect().Model(&events).
		Where("domain = ?", domain).
		Where("bucket_start BETWEEN ? AND ?", start, end).
		Where("is_anomaly = ?", true).
		Order("source_ip asc").
		Order("bucket_start asc").
		Order("score desc").
		Order("total_bytes desc").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return events, nil
}
