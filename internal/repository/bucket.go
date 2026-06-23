package repository

import (
	"context"
	"errors"
	"time"
	"traffic-guarder/internal/model"

	"github.com/uptrace/bun"
)

type BucketRepository interface {
	UpsertBucket(ctx context.Context, bucket *model.TrafficBucket) (*model.TrafficBucket, error)
	BucketByDomainAndStart(ctx context.Context, domain string, startTime time.Time) ([]*model.TrafficBucket, error)
}

type bucketRepository struct {
	db *bun.DB
}

func NewBucketRepository(db *bun.DB) BucketRepository {
	return &bucketRepository{db: db}
}

func (r *bucketRepository) UpsertBucket(ctx context.Context, bucket *model.TrafficBucket) (*model.TrafficBucket, error) {
	_, err := r.db.NewInsert().
		Model(bucket).
		On("CONFLICT (bucket_start, domain, source_ip) DO UPDATE").
		Set("request_count = traffic_buckets.request_count + EXCLUDED.request_count").
		Set("request_bytes_sum = traffic_buckets.request_bytes_sum + EXCLUDED.request_bytes_sum").
		Set("response_bytes_sum = traffic_buckets.response_bytes_sum + EXCLUDED.response_bytes_sum").
		Set("total_bytes_sum = traffic_buckets.total_bytes_sum + EXCLUDED.total_bytes_sum").
		Set("nx_domain_count = traffic_buckets.nx_domain_count + EXCLUDED.nx_domain_count").
		Set("servfail_count = traffic_buckets.servfail_count + EXCLUDED.servfail_count").
		Set("no_error_count = traffic_buckets.no_error_count + EXCLUDED.no_error_count").
		Set("latency_sum_ms = traffic_buckets.latency_sum_ms + EXCLUDED.latency_sum_ms").
		Set("country = EXCLUDED.country").
		Set("asn = EXCLUDED.asn").
		Set("protocol = EXCLUDED.protocol").
		Set("updated_at = NOW()").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (r *bucketRepository) BucketByDomainAndStart(ctx context.Context, domain string, startTime time.Time) ([]*model.TrafficBucket, error) {
	var bucket []*model.TrafficBucket
	err := r.db.NewSelect().Model(&bucket).
		Where("domain = ?", domain).
		Where("bucket_start = ?", startTime).
		Order("total_bytes_sum DESC").
		Scan(ctx)
	if err != nil {
		return nil, errors.New("select buckets from domain error:" + err.Error())
	}
	return bucket, nil
}
