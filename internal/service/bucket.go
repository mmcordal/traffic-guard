package service

import (
	"context"
	"errors"
	"fmt"
	"traffic-guarder/internal/infrastructure/cache"
	"traffic-guarder/internal/infrastructure/config"
	"traffic-guarder/internal/model"
	"traffic-guarder/internal/repository"
)

type BucketService interface {
	UpsertBucket(ctx context.Context, log *model.TrafficLog) error
}

type bucketService struct {
	br  repository.BucketRepository
	bc  cache.BucketCache
	cfg config.AnalyzeConfig
}

func NewBucketService(br repository.BucketRepository, bc cache.BucketCache, cfg config.AnalyzeConfig) BucketService {
	return &bucketService{br: br, bc: bc, cfg: cfg}
}

func (s *bucketService) UpsertBucket(ctx context.Context, log *model.TrafficLog) error {
	bucketStart := log.Timestamp.Truncate(s.cfg.BucketWindow())

	nxDomainCount := int64(0)
	servfailCount := int64(0)
	noErrorCount := int64(0)

	switch log.ResponseCode {
	case model.ResponseCodeNXDomain:
		nxDomainCount = 1
	case model.ResponseCodeServfail:
		servfailCount = 1
	case model.ResponseCodeNoError:
		noErrorCount = 1
	}

	bucket := &model.TrafficBucket{
		BucketStart:      bucketStart,
		Domain:           log.Domain,
		SourceIP:         log.SourceIP,
		RequestCount:     1,
		RequestBytesSum:  log.RequestBytes,
		ResponseBytesSum: log.ResponseBytes,
		TotalBytesSum:    log.TotalBytes,
		NXDomainCount:    nxDomainCount,
		ServfailCount:    servfailCount,
		NoErrorCount:     noErrorCount,
		LatencySumMs:     log.LatencyMs,
		Country:          log.Country,
		ASN:              log.ASN,
		Protocol:         log.Protocol,
	}

	_, err := s.br.UpsertBucket(ctx, bucket)
	if err != nil {
		return fmt.Errorf("upsert bucket error: %v", err)
	}

	if err := s.bc.IncrementDomainBucketFromLog(ctx, log); err != nil {
		return errors.New("increment domain bucket from bc error:" + err.Error())
	}

	return nil
}
