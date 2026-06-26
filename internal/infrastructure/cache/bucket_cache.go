package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
	"traffic-guarder/internal/infrastructure/config"
	"traffic-guarder/internal/model"

	"github.com/redis/go-redis/v9"
)

type BucketCache interface {
	IncrementDomainBucketFromLog(ctx context.Context, log *model.TrafficLog) error
	GetDomainsByBucket(ctx context.Context, bucketStart time.Time) ([]string, error)
	GetDomainBucket(ctx context.Context, domain string, bucketStart time.Time) (map[string]string, error)
	GetPreviousBucketMinutes(ctx context.Context, domain string, before time.Time, limit int64) ([]string, error)
}

type bucketCache struct {
	cache *RedisClient
	cfg   config.AnalyzeConfig
}

func NewBucketCache(c *RedisClient, cfg config.AnalyzeConfig) BucketCache {
	return &bucketCache{cache: c, cfg: cfg}
}

func (c *bucketCache) IncrementDomainBucketFromLog(ctx context.Context, log *model.TrafficLog) error {
	bucketStart := log.Timestamp.Truncate(c.cfg.BucketWindow())
	bucketUnix := bucketStart.Unix()

	key := fmt.Sprintf("tg:domain_bucket:%s:%d", log.Domain, bucketUnix)
	minutesKey := fmt.Sprintf("tg:domain_minutes:%s", log.Domain)
	bucketDomainsKey := fmt.Sprintf("tg:bucket_domains:%d", bucketUnix)

	if c.cache.Client == nil {
		return errors.New("redis client not initialized")
	}

	if err := c.cache.HIncrBy(ctx, key, "request_count", 1); err != nil {
		return errors.New("redis hIncrBy bucket_cache --> request_count error: " + err.Error())
	}

	if err := c.cache.HIncrBy(ctx, key, "request_bytes_sum", log.RequestBytes); err != nil {
		return errors.New("redis hIncrBy bucket_cache --> request_sum error: " + err.Error())
	}

	if err := c.cache.HIncrBy(ctx, key, "response_bytes_sum", log.ResponseBytes); err != nil {
		return errors.New("redis hIncrBy bucket_cache --> response_sum error: " + err.Error())
	}

	if err := c.cache.HIncrBy(ctx, key, "total_bytes_sum", log.TotalBytes); err != nil {
		return errors.New("redis hIncrBy bucket_cache --> total_bytes_sum error: " + err.Error())
	}

	if model.ResponseCodeNXDomain == log.ResponseCode {
		if err := c.cache.HIncrBy(ctx, key, "nx_domain_count", 1); err != nil {
			return errors.New("redis hIncrBy bucket_cache --> nx_domain_count error: " + err.Error())
		}
	}

	if model.ResponseCodeServfail == log.ResponseCode {
		if err := c.cache.HIncrBy(ctx, key, "servfail_count", 1); err != nil {
			return errors.New("redis hIncrBy bucket_cache --> servfail_count error: " + err.Error())
		}
	}

	if model.ResponseCodeNoError == log.ResponseCode {
		if err := c.cache.HIncrBy(ctx, key, "no_error_count", 1); err != nil {
			return errors.New("redis hIncrBy bucket_cache --> no_error_count error: " + err.Error())
		}
	}

	if err := c.cache.HIncrBy(ctx, key, "latency_sum_ms", log.LatencyMs); err != nil {
		return errors.New("redis hIncrBy bucket_cache --> latency_sum_ms error: " + err.Error())
	}

	if err := c.cache.HSet(ctx, key,
		"bucket_start", bucketStart.Format(time.RFC3339),
		"domain", log.Domain,
	); err != nil {
		return errors.New("redis hSet bucket_cache error: " + err.Error())
	}

	if err := c.cache.ZAdd(ctx, minutesKey, redis.Z{
		Score:  float64(bucketUnix),
		Member: bucketUnix,
	}); err != nil {
		return errors.New("redis zAdd bucket_cache error: " + err.Error())
	}

	if err := c.cache.SAdd(ctx, bucketDomainsKey, log.Domain); err != nil {
		return errors.New("redis sAdd bucket_cache --> bucket_domains error: " + err.Error())
	}

	_ = c.cache.Expire(ctx, key, c.cfg.BucketTTL())
	_ = c.cache.Expire(ctx, minutesKey, c.cfg.BucketTTL())
	_ = c.cache.Expire(ctx, bucketDomainsKey, c.cfg.BucketTTL())

	return nil
}

func (c *bucketCache) GetDomainsByBucket(ctx context.Context, bucketStart time.Time) ([]string, error) {
	key := fmt.Sprintf("tg:bucket_domains:%d", bucketStart.Unix())

	return c.cache.SMembers(ctx, key)

}

func (c *bucketCache) GetDomainBucket(ctx context.Context, domain string, bucketStart time.Time) (map[string]string, error) {
	key := fmt.Sprintf("tg:domain_bucket:%s:%d", domain, bucketStart.Unix())

	return c.cache.HGetAll(ctx, key)
}

func (c *bucketCache) GetPreviousBucketMinutes(ctx context.Context, domain string, before time.Time, limit int64) ([]string, error) {
	key := fmt.Sprintf("tg:domain_minutes:%s", domain)

	return c.cache.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Max:    fmt.Sprintf("%d", before.Unix()),
		Min:    "-inf",
		Offset: 0,
		Count:  limit,
	})
}
