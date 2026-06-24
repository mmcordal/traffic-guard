package service

import (
	"math"
	"time"
	"traffic-guarder/internal/model"
	"traffic-guarder/internal/viewmodel"
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

func help(anomalyEvents []*model.AnomalyEvent, domain string) *viewmodel.ExclusionResponse {
	vm := &viewmodel.ExclusionResponse{
		Domain:        domain,
		AnomalyEvents: []*viewmodel.AnomalyEvent{},
		DangerousIPs:  []string{},
	}
	if len(anomalyEvents) == 0 {
		return vm
	}

	currentIP := anomalyEvents[0].SourceIP
	vm.DangerousIPs = append(vm.DangerousIPs, currentIP)

	current := newAnomalyEventVM(anomalyEvents[0], domain)

	scoreSum := anomalyEvents[0].Score
	scoreCount := int64(1)

	reasonBytes := make(map[model.AnomalyReason]int64)
	reasonBytes[anomalyEvents[0].Reason] = anomalyEvents[0].TotalBytes

	for _, event := range anomalyEvents[1:] {
		if event.SourceIP != currentIP {
			current.Score = scoreSum / float64(scoreCount)
			current.MostReason = string(mostReasonByBytes(reasonBytes))
			vm.AnomalyEvents = append(vm.AnomalyEvents, current)

			currentIP = event.SourceIP
			vm.DangerousIPs = append(vm.DangerousIPs, currentIP)

			current = newAnomalyEventVM(event, domain)

			scoreSum = event.Score
			scoreCount = 1

			reasonBytes = make(map[model.AnomalyReason]int64)
			reasonBytes[event.Reason] = event.TotalBytes

			continue
		}

		current.AttackEnded = event.BucketStart.Add(time.Minute).Format(time.RFC3339)

		current.TotalBytes += event.TotalBytes
		current.TotalRequests += event.RequestCount
		current.TotalNXDomain += event.NXDomainCount
		current.TotalServfail += event.ServfailCount
		current.TotalNoError += event.NoErrorCount

		scoreSum += event.Score
		scoreCount++

		reasonBytes[event.Reason] += event.TotalBytes

	}
	current.Score = scoreSum / float64(scoreCount)
	current.MostReason = string(mostReasonByBytes(reasonBytes))
	vm.AnomalyEvents = append(vm.AnomalyEvents, current)

	fillExclusionTotals(vm)

	return vm
}

func newAnomalyEventVM(event *model.AnomalyEvent, domain string) *viewmodel.AnomalyEvent {
	return &viewmodel.AnomalyEvent{
		AttackStarted: event.BucketStart.Format(time.RFC3339),
		AttackEnded:   event.BucketStart.Add(time.Minute).Format(time.RFC3339),

		Domain: domain,
		IP:     event.SourceIP,

		Score:      event.Score,
		MostReason: string(event.Reason),

		TotalBytes:    event.TotalBytes,
		TotalRequests: event.RequestCount,

		TotalNXDomain: event.NXDomainCount,
		TotalServfail: event.ServfailCount,
		TotalNoError:  event.NoErrorCount,

		Protocol: string(event.Protocol),
		Country:  (event.Country),
		ASN:      event.ASN,
	}
}

func mostReasonByBytes(values map[model.AnomalyReason]int64) model.AnomalyReason {
	maxReason := model.ReasonBytesSpike
	maxVal := int64(0)

	for reason, value := range values {
		if value > maxVal {
			maxReason = reason
			maxVal = value
		}
	}
	return maxReason
}

func fillExclusionTotals(vm *viewmodel.ExclusionResponse) {
	scoreSum := float64(0)
	scoreCount := float64(0)

	asnBytes := make(map[string]int64)
	countryBytes := make(map[string]int64)
	protocolBytes := make(map[string]int64)
	reasonBytes := make(map[string]int64)

	for _, event := range vm.AnomalyEvents {
		scoreSum += event.Score
		scoreCount++

		vm.SpentTotalBytes += event.TotalBytes
		vm.TotalRequests += event.TotalRequests

		vm.TotalNXDomain += event.TotalNXDomain
		vm.TotalServfail += event.TotalServfail
		vm.TotalNoError += event.TotalNoError

		asnBytes[event.ASN] += event.TotalBytes
		countryBytes[event.Country] += event.TotalBytes
		protocolBytes[event.Protocol] += event.TotalBytes
		reasonBytes[string(event.MostReason)] += event.TotalBytes
	}
	if scoreCount > 0 {
		vm.AverageScore = scoreSum / scoreCount
	}

	vm.MostASN = maxStringsByBytes(asnBytes)
	vm.MostCountry = maxStringsByBytes(countryBytes)
	vm.MostProtocol = maxStringsByBytes(protocolBytes)
	vm.MostAnomalyReason = maxStringsByBytes(reasonBytes)
}

func maxStringsByBytes(values map[string]int64) string {
	maxKey := ""
	maxVal := int64(0)

	for k, v := range values {
		if k == "" {
			continue
		}

		if v > maxVal {
			maxKey = k
			maxVal = v
		}
	}
	return maxKey
}
