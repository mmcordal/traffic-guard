package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TrafficBucket struct {
	bun.BaseModel `bun:"table:traffic_buckets,alias:tb"`
	CoreModel

	BucketStart time.Time `bun:",notnull,unique:idx_bucket_domain_ip"`

	Domain   string `bun:",type:varchar(255),notnull,unique:idx_bucket_domain_ip"`
	SourceIP string `bun:",type:varchar(45),notnull,unique:idx_bucket_domain_ip"`

	RequestCount int64 `bun:",notnull,check:request_count >= 0"`

	RequestBytesSum  int64 `bun:",notnull,check:request_bytes_sum >= 0"`
	ResponseBytesSum int64 `bun:",notnull,check:response_bytes_sum >= 0"`
	TotalBytesSum    int64 `bun:",notnull,check:total_bytes_sum >= 0"`

	NXDomainCount int64 `bun:",notnull,check:nx_domain_count >= 0"`
	ServfailCount int64 `bun:",notnull,check:servfail_count >= 0"`
	NoErrorCount  int64 `bun:",notnull,check:no_error_count >= 0"`

	LatencySumMs int64 `bun:",notnull,check:latency_sum_ms >= 0"`

	Country  string   `bun:",type:varchar(2),nullzero"`
	ASN      string   `bun:",type:varchar(32),nullzero"`
	Protocol Protocol `bun:",type:varchar(10),notnull"`
}
