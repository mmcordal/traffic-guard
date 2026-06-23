package repository

import (
	"context"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type DomainCheck interface {
	Create(ctx context.Context, check *model.DomainAnomalyCheck) error
}

type domainCheck struct {
	db *bun.DB
}

func NewDomainCheck(db *bun.DB) DomainCheck {
	return &domainCheck{db: db}
}

func (d *domainCheck) Create(ctx context.Context, check *model.DomainAnomalyCheck) error {
	_, err := d.db.NewInsert().Model(check).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
