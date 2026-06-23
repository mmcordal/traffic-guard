package model

import (
	"time"

	"github.com/uptrace/bun"
)

type DomainAnomalyCheck struct {
	bun.BaseModel `bun:"table:domain_anomaly_checks,alias:dac"`
	CoreModel

	BucketStart time.Time `bun:",notnull"`
	Domain      string    `bun:",type:varchar(255),notnull"`

	Score     float64 `bun:",notnull"`
	IsAnomaly bool    `bun:",notnull,default:false"`
	Severity  string  `bun:",type:varchar(30),notnull"`

	Reason AnomalyReason `bun:",type:varchar(80),notnull"`

	TotalBytes   int64 `bun:",notnull,check:total_bytes >= 0"`
	RequestCount int64 `bun:",notnull,check:request_count >= 0"`

	BytesScore    float64 `bun:",notnull,default:0"`
	RequestScore  float64 `bun:",notnull,default:0"`
	NXScore       float64 `bun:",notnull,default:0"`
	ServfailScore float64 `bun:",notnull,default:0"`
}
