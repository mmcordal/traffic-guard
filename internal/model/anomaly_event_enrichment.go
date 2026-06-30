package model

type AnomalyEventEnrichment struct {
	QueryType    QueryType
	LatencySumMs int64
	AvgLatencyMs float64
}
