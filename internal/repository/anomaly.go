package repository

import (
	"context"
	"time"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type AnomalyRepository interface {
	CreateAnomalyEvent(ctx context.Context, event *model.AnomalyEvent) error
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

func (r *anomalyRepository) GetEventsWithStartAndEndTime(ctx context.Context, domain string, start, end time.Time) ([]*model.AnomalyEvent, error) {
	var events []*model.AnomalyEvent
	err := r.db.NewSelect().Model(&events).
		Where("domain = ?", domain).
		Where("bucket_start BETWEEN ? AND ?", start, end).
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
