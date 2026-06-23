package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TrafficLog struct {
	bun.BaseModel `bun:"table:traffic_logs,alias:tl"`
	CoreModel

	Timestamp time.Time `bun:",notnull"`

	Domain    string `bun:",type:varchar(255),notnull"`
	QueryName string `bun:",type:varchar(255),notnull"`
	SourceIP  string `bun:",type:varchar(45),notnull"`

	QueryType    QueryType    `bun:",type:varchar(20),notnull"`
	ResponseCode ResponseCode `bun:",type:varchar(20),notnull"`

	RequestBytes  int64 `bun:",notnull,check:request_bytes >= 0"`
	ResponseBytes int64 `bun:",notnull,check:response_bytes >= 0"`
	TotalBytes    int64 `bun:",notnull,check:total_bytes >= 0"`

	Protocol  Protocol `bun:",type:varchar(10),notnull"`
	LatencyMs int64    `bun:",notnull,check:latency_ms >= 0"`

	Country string `bun:",type:varchar(2),nullzero"`  // TR, US, DE gibi
	ASN     string `bun:",type:varchar(32),nullzero"` // AS9121 gibi
}
