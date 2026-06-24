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

func mostReason(bytes, req, nx, serv, noerr int64) model.AnomalyReason {
	maxVal := bytes
	maxReason := model.ReasonBytesSpike

	if req > maxVal {
		maxReason = model.ReasonRequestSpike
	}
	if nx > maxVal {
		maxReason = model.ReasonNXDomainSpike
	}
	if serv > maxVal {
		maxReason = model.ReasonServfailSpike
	}
	if noerr > maxVal {
		maxReason = model.ReasonNoErrorSpike
	}
	return maxReason
}

func help(anomalyEvents []*model.AnomalyEvent, domain string) *viewmodel.ExclusionResponse {
	ips := []string{}
	ips = append(ips, anomalyEvents[0].SourceIP)
	aes := []*viewmodel.AnomalyEvent{}
	ae := &viewmodel.AnomalyEvent{
		AttackStarted: anomalyEvents[0].BucketStart.Format(time.RFC3339),
		AttackEnded:   anomalyEvents[0].BucketStart.Add(time.Minute).Format(time.RFC3339),
		Domain:        domain,
		IP:            anomalyEvents[0].SourceIP,
		Score:         anomalyEvents[0].Score,
		MostReason:    string(anomalyEvents[0].Reason),
		TotalBytes:    anomalyEvents[0].TotalBytes,
		TotalRequests: anomalyEvents[0].RequestCount,
		TotalNXDomain: anomalyEvents[0].NXDomainCount,
		TotalServfail: anomalyEvents[0].ServfailCount,
		TotalNoError:  anomalyEvents[0].NoErrorCount,
		Protocol:      string(anomalyEvents[0].Protocol),
		Country:       anomalyEvents[0].Country,
		ASN:           anomalyEvents[0].ASN,
	}
	reqSpike := int64(0)
	bytesSpike := int64(0)
	nxSpike := int64(0)
	servfailSpike := int64(0)
	noErrorSpike := int64(0)

	sayac := 1.0

	ip := ""
	for i, event := range anomalyEvents {
		if i == 0 { // İLKİM
			continue
		}
		if i == len(anomalyEvents)-1 && event.SourceIP == ip { // en sonda ve bi öncekiyle aynı yani son
			ae.AttackEnded = event.BucketStart.Add(time.Minute).Format(time.RFC3339)
			sayac++
			ae.Score = (ae.Score + event.Score) / sayac
			ae.TotalBytes += event.TotalBytes
			ae.TotalRequests += event.RequestCount
			ae.TotalNXDomain += event.NXDomainCount
			ae.TotalServfail += event.ServfailCount
			ae.TotalNoError += event.NoErrorCount
			switch event.Reason {
			case model.ReasonRequestSpike:
				reqSpike++
			case model.ReasonBytesSpike:
				bytesSpike++
			case model.ReasonNXDomainSpike:
				nxSpike++
			case model.ReasonServfailSpike:
				servfailSpike++
			case model.ReasonNoErrorSpike:
				noErrorSpike++
			}
			ae.MostReason = string(mostReason(bytesSpike, reqSpike, nxSpike, servfailSpike, noErrorSpike))
			aes = append(aes, ae)
		} else { // tekim ve en sonum
			ips = append(ips, event.SourceIP)
			ae.AttackStarted = event.BucketStart.Format(time.RFC3339)
			ae.AttackEnded = event.BucketStart.Add(time.Minute).Format(time.RFC3339)
			ae.Domain = domain
			ae.IP = event.SourceIP
			ae.Score = event.Score
			ae.MostReason = string(event.Reason)
			ae.TotalBytes = event.TotalBytes
			ae.TotalRequests = event.RequestCount
			ae.TotalNXDomain = event.NXDomainCount
			ae.TotalServfail = event.ServfailCount
			ae.TotalNoError = event.NoErrorCount
			ae.Protocol = string(event.Protocol)
			ae.Country = event.Country
			ae.ASN = event.ASN
			aes = append(aes, ae)
		}
		// iterator.ip == next.ip
		if event.SourceIP == anomalyEvents[i+1].SourceIP {
			// iterator.ip != prev.ip --> iterator.ip == next.ip && iterator.ip != prev.ip
			if event.SourceIP != ip { // BAŞTAYIM
				ip = event.SourceIP
				ips = append(ips, ip)
				sayac = 1.0
				ae = &viewmodel.AnomalyEvent{
					AttackStarted: event.BucketStart.Format(time.RFC3339),

					Domain: domain,
					IP:     ip,

					Score: event.Score,

					TotalBytes:    event.TotalBytes,
					TotalRequests: event.RequestCount,

					TotalNXDomain: event.NXDomainCount,
					TotalServfail: event.ServfailCount,
					TotalNoError:  event.NoErrorCount,

					Protocol: string(event.Protocol),
					Country:  event.Country,
					ASN:      event.ASN,
				}
				switch event.Reason {
				case model.ReasonRequestSpike:
					reqSpike = 1
					bytesSpike = 0
					nxSpike = 0
					servfailSpike = 0
					noErrorSpike = 0
				case model.ReasonBytesSpike:
					bytesSpike = 1
					reqSpike = 0
					nxSpike = 0
					servfailSpike = 0
					noErrorSpike = 0
				case model.ReasonNXDomainSpike:
					nxSpike = 1
					reqSpike = 0
					noErrorSpike = 0
					servfailSpike = 0
					bytesSpike = 0
				case model.ReasonServfailSpike:
					servfailSpike = 1
					bytesSpike = 0
					reqSpike = 0
					noErrorSpike = 0
					servfailSpike = 0
				case model.ReasonNoErrorSpike:
					noErrorSpike = 1
					bytesSpike = 0
					reqSpike = 0
					noErrorSpike = 0
					servfailSpike = 0
				}

			} else { // iterator.ip == prev.ip --> iterator.ip == next.ip && iterator.ip == prev.ip
				// ORTADAYIM
				sayac++
				ae.Score += event.Score
				ae.TotalBytes += event.TotalBytes
				ae.TotalRequests += event.RequestCount
				ae.TotalNXDomain += event.NXDomainCount
				ae.TotalServfail += event.ServfailCount
				ae.TotalNoError += event.NoErrorCount

				switch event.Reason {
				case model.ReasonRequestSpike:
					reqSpike++
				case model.ReasonBytesSpike:
					bytesSpike++
				case model.ReasonNXDomainSpike:
					nxSpike++
				case model.ReasonServfailSpike:
					servfailSpike++
				case model.ReasonNoErrorSpike:
					noErrorSpike++
				}
			}

		} else { // iterator.ip != next.ip
			if event.SourceIP == ip { // iterator.ip == prev.ip --> iterator.ip != next.ip && iterator.ip == prev.ip
				// SONDAYIM
				ae.AttackEnded = event.BucketStart.Add(time.Minute).Format(time.RFC3339)
				sayac++
				ae.Score = (ae.Score + event.Score) / sayac
				ae.TotalBytes += event.TotalBytes
				ae.TotalRequests += event.RequestCount
				ae.TotalNXDomain += event.NXDomainCount
				ae.TotalServfail += event.ServfailCount
				ae.TotalNoError += event.NoErrorCount
				switch event.Reason {
				case model.ReasonRequestSpike:
					reqSpike++
				case model.ReasonBytesSpike:
					bytesSpike++
				case model.ReasonNXDomainSpike:
					nxSpike++
				case model.ReasonServfailSpike:
					servfailSpike++
				case model.ReasonNoErrorSpike:
					noErrorSpike++
				}
				ae.MostReason = string(mostReason(bytesSpike, reqSpike, nxSpike, servfailSpike, noErrorSpike))
				aes = append(aes, ae)
			} else { // iterator.ip != prev.ip -->iterator.ip != next.ip && iterator.ip != prev.ip
				// TEKİM
				ip = event.SourceIP
				ips = append(ips, ip)
				ae.AttackStarted = event.BucketStart.Format(time.RFC3339)
				ae.AttackEnded = event.BucketStart.Add(time.Minute).Format(time.RFC3339)
				ae.Domain = domain
				ae.IP = event.SourceIP
				ae.Score = event.Score
				ae.MostReason = string(event.Reason)
				ae.TotalBytes = event.TotalBytes
				ae.TotalRequests = event.RequestCount
				ae.TotalNXDomain = event.NXDomainCount
				ae.TotalServfail = event.ServfailCount
				ae.TotalNoError = event.NoErrorCount
				ae.Protocol = string(event.Protocol)
				ae.Country = event.Country
				ae.ASN = event.ASN
				aes = append(aes, ae)
				sayac = 0.0
			}
		}
	}

	vm := mosts(aes)
	vm.Domain = domain
	vm.AnomalyEvents = aes
	vm.DangerousIPs = ips

	score := 0.0
	count := 0.0

	for _, v := range aes {
		vm.SpentTotalBytes += v.TotalBytes
		vm.TotalRequests += v.TotalRequests
		vm.TotalNXDomain += v.TotalNXDomain
		vm.TotalServfail += v.TotalServfail
		vm.TotalNoError += v.TotalNoError
		score += v.Score
		count++
	}
	vm.AverageScore = score / count

	return vm
}

func mosts(event []*viewmodel.AnomalyEvent) *viewmodel.ExclusionResponse {
	var asnMap map[string]int64
	var countryMap map[string]int64
	var protocolMap map[string]int64
	var reasonMap map[string]int64

	for _, v := range event {
		asnMap[v.ASN] += 1
		countryMap[v.Country] += 1
		protocolMap[string(v.Protocol)] += 1
		reasonMap[string(v.MostReason)] += 1
	}

	vm := new(viewmodel.ExclusionResponse)

	maxK := ""
	maxV := int64(0)
	for k, v := range asnMap {
		if v > maxV {
			maxK = k
			maxV = v
		}
	}
	vm.MostASN = maxK

	maxK = ""
	maxV = int64(0)
	for k, v := range countryMap {
		if v > maxV {
			maxK = k
			maxV = v
		}
	}
	vm.MostCountry = maxK

	maxK = ""
	maxV = int64(0)
	for k, v := range protocolMap {
		if v > maxV {
			maxK = k
			maxV = v
		}
	}
	vm.MostProtocol = maxK

	maxK = ""
	maxV = int64(0)
	for k, v := range reasonMap {
		if v > maxV {
			maxK = k
			maxV = v
		}
	}
	vm.MostAnomalyReason = maxK
	return vm
}
