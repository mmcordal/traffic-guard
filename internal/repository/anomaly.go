package repository

import (
	"context"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type AnomalyRepository interface {
	CreateAnomalyEvent(ctx context.Context, event *model.AnomalyEvent) error
}

type analyzeRepository struct {
	db *bun.DB
}

func NewAnomalyRepository(db *bun.DB) AnomalyRepository {
	return &analyzeRepository{db: db}
}

func (r *analyzeRepository) CreateAnomalyEvent(ctx context.Context, event *model.AnomalyEvent) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	return err
}
