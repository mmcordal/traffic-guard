package model

import (
	"time"

	"github.com/uptrace/bun"
)

type AnomalyEvent struct {
	bun.BaseModel `bun:"table:anomaly_events,alias:ae"`
	CoreModel

	BucketStart time.Time `bun:",notnull"`

	Domain   string `bun:",type:varchar(255),notnull"`
	SourceIP string `bun:",type:varchar(45),notnull"`

	Score     float64       `bun:",notnull"`
	IsAnomaly bool          `bun:",notnull,default:false"`
	Reason    AnomalyReason `bun:",type:varchar(80),notnull"`

	TotalBytes   int64 `bun:",notnull,check:total_bytes >= 0"`
	RequestCount int64 `bun:",notnull,check:request_count >= 0"`

	NXDomainCount int64 `bun:",notnull,check:request_count >= 0"`
	ServfailCount int64 `bun:",notnull,check:servfail >= 0"`
	NoErrorCount  int64 `bun:",notnull,check:no_error >= 0"`

	QueryType QueryType `bun:",type:varchar(20),nullzero"`
	Protocol  Protocol  `bun:",type:varchar(10),nullzero"`

	Country string `bun:",type:varchar(2),nullzero"`
	ASN     string `bun:",type:varchar(32),nullzero"`
}
