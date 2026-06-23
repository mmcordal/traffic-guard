package service

import (
	"math"
	"traffic-guarder/internal/model"
)

const (
	WarningThreshold  = 2.0
	AnomalyThreshold  = 3.0
	CriticalThreshold = 5.0

	MinBaselineBuckets = 5

	MinRequestCount = 30
	MinTotalBytes   = 10 * 1024 // 10 KB
)

type AnomalySeverity string

const (
	SeverityNormal   AnomalySeverity = "normal"
	SeverityWarning  AnomalySeverity = "warning"
	SeverityHigh     AnomalySeverity = "high"
	SeverityCritical AnomalySeverity = "critical"
)

func classifyAnomaly(score, currentReq, currentBytes float64) (bool, AnomalySeverity) {
	if currentReq < MinRequestCount && currentBytes < MinTotalBytes {
		return false, SeverityNormal // bytes < 10KB ve request < 30 ise
	}

	switch {
	case score >= CriticalThreshold:
		return true, SeverityCritical // 5 < x
	case score >= AnomalyThreshold:
		return true, SeverityHigh // 3 < x < 5
	case score >= WarningThreshold:
		return false, SeverityWarning // 2 < x < 3
	default:
		return false, SeverityNormal // x < 2
	}
}

func zScore(current float64, prev []float64) float64 {
	if len(prev) < 2 {
		return 0
	}
	avg := mean(prev)
	s := stddev(prev, avg)

	if s == 0 {
		if current > avg {
			return current - avg
		}
		return 0
	}
	return (current - avg) / s
}

func mean(array []float64) float64 {
	sum := 0.0
	for _, v := range array {
		sum += v
	}
	return sum / float64(len(array))
}

func maxScore(bytesScore, reqScore, nxScore, servfailScore, noErrorScore float64) (float64, string) {
	scores := []struct {
		value  float64
		reason string
	}{
		{value: bytesScore, reason: "bytes_spike"},
		{value: reqScore, reason: "request_spike"},
		{value: nxScore, reason: "nxdomain_spike"},
		{value: servfailScore, reason: "servfail_spike"},
		{value: noErrorScore, reason: "noerror_spike"},
	}

	maxValue := scores[0].value
	maxReason := scores[0].reason

	for _, score := range scores[1:] {
		if score.value > maxValue {
			maxValue = score.value
			maxReason = score.reason
		}
	}

	return maxValue, maxReason
}

func stddev(array []float64, avg float64) float64 {
	sum := 0.0
	for _, v := range array {
		diff := v - avg
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(array)-1))
}

func forRequestSpike(buckets []*model.TrafficBucket) *model.TrafficBucket {
	dangerousBucket := new(model.TrafficBucket)
	var maxCount int64
	for _, bucket := range buckets {
		if bucket.RequestCount > maxCount {
			maxCount = bucket.RequestCount
			dangerousBucket = bucket
		}
	}
	return dangerousBucket
}

func forBytesSpike(buckets []*model.TrafficBucket) *model.TrafficBucket {
	dangerousBucket := new(model.TrafficBucket)
	var maxCount int64
	for _, bucket := range buckets {
		if bucket.TotalBytesSum > maxCount {
			maxCount = bucket.TotalBytesSum
			dangerousBucket = bucket
		}
	}
	return dangerousBucket
}

func forNXSpike(buckets []*model.TrafficBucket) *model.TrafficBucket {
	dangerousBucket := new(model.TrafficBucket)
	var maxCount int64
	for _, bucket := range buckets {
		if bucket.NXDomainCount > maxCount {
			maxCount = bucket.NXDomainCount
			dangerousBucket = bucket
		}
	}
	return dangerousBucket
}

func forServerfailSpike(buckets []*model.TrafficBucket) *model.TrafficBucket {
	dangerousBucket := new(model.TrafficBucket)
	var maxCount int64
	for _, bucket := range buckets {
		if bucket.ServfailCount > maxCount {
			maxCount = bucket.ServfailCount
			dangerousBucket = bucket
		}
	}
	return dangerousBucket
}

func forNoErrorSpike(buckets []*model.TrafficBucket) *model.TrafficBucket {
	dangerousBucket := new(model.TrafficBucket)
	var maxCount int64
	for _, bucket := range buckets {
		if bucket.NoErrorCount > maxCount {
			maxCount = bucket.NoErrorCount
			dangerousBucket = bucket
		}
	}
	return dangerousBucket
}
