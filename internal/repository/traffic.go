package repository

import (
	"context"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type TrafficRepository interface {
	CreateLog(ctx context.Context, log *model.TrafficLog) (*model.TrafficLog, error)
}

type trafficRepository struct {
	db *bun.DB
}

func NewTrafficRepository(db *bun.DB) TrafficRepository {
	return &trafficRepository{db: db}
}

func (r *trafficRepository) CreateLog(ctx context.Context, log *model.TrafficLog) (*model.TrafficLog, error) {
	_, err := r.db.NewInsert().Model(log).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return log, nil
}
