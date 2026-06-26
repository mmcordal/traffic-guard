package service

import (
	"context"
	"errors"
	"strconv"
	"time"
	"traffic-guarder/internal/infrastructure/cache"
	"traffic-guarder/internal/infrastructure/config"
	"traffic-guarder/internal/model"
	"traffic-guarder/internal/repository"
	"traffic-guarder/internal/viewmodel"
)

type AnomalyService interface {
	AnalyzeCompletedBucket(ctx context.Context, bucketStart time.Time) error
	GetAnomalyEvents(ctx context.Context, vm *viewmodel.ExclusionRequest) (*viewmodel.ExclusionResponse, error)
}

type anomalyService struct {
	ar  repository.AnomalyRepository
	br  repository.BucketRepository
	dc  repository.DomainCheck
	bc  cache.BucketCache
	cfg config.AnalyzeConfig
}

func NewAnomalyService(ar repository.AnomalyRepository,
	dc repository.DomainCheck,
	bc cache.BucketCache,
	br repository.BucketRepository,
	cfg config.AnalyzeConfig,
) AnomalyService {
	return &anomalyService{ar: ar, dc: dc, bc: bc, br: br, cfg: cfg}
}

func (s *anomalyService) AnalyzeCompletedBucket(ctx context.Context, bucketStart time.Time) error {
	domains, err := s.bc.GetDomainsByBucket(ctx, bucketStart)
	if err != nil {
		return errors.New("anomalyService GetDomainsByBucket error: " + err.Error())
	}

	for _, domain := range domains {
		currentMap, err := s.bc.GetDomainBucket(ctx, domain, bucketStart)
		if err != nil {
			return errors.New("anomalyService GetDomainBucket error: " + err.Error())
		}

		currentReq, _ := strconv.ParseFloat(currentMap["request_count"], 64)
		currentBytes, _ := strconv.ParseFloat(currentMap["total_bytes_sum"], 64)
		currentNX, _ := strconv.ParseFloat(currentMap["nx_domain_count"], 64)
		currentServfail, _ := strconv.ParseFloat(currentMap["servfail_count"], 64)
		currentNoError, _ := strconv.ParseFloat(currentMap["no_error_count"], 64)

		previousMinutes, err := s.bc.GetPreviousBucketMinutes(
			ctx,
			domain,
			bucketStart.Add(-s.cfg.BucketWindow()),
			s.cfg.HistoryLimit(),
		)
		if err != nil {
			return errors.New("anomalyService GetPreviousBucketMinutes error: " + err.Error())
		}

		var prevReq []float64
		var prevBytes []float64
		var prevNX []float64
		var prevServfail []float64
		var prevNoError []float64

		for _, minuteStr := range previousMinutes {
			minuteUnix, err := strconv.ParseInt(minuteStr, 10, 64)
			if err != nil {
				continue
			}

			prevBucketStart := time.Unix(minuteUnix, 0)

			prevMap, err := s.bc.GetDomainBucket(ctx, domain, prevBucketStart)
			if err != nil {
				return errors.New("anomalyService GetDomainBucket error2: " + err.Error())
			}

			bytesVal, _ := strconv.ParseFloat(prevMap["total_bytes_sum"], 64)
			reqVal, _ := strconv.ParseFloat(prevMap["request_count"], 64)
			nxVal, _ := strconv.ParseFloat(prevMap["nx_domain_count"], 64)
			servfailVal, _ := strconv.ParseFloat(prevMap["servfail_count"], 64)
			noErrorVal, _ := strconv.ParseFloat(prevMap["no_error_count"], 64)

			prevBytes = append(prevBytes, bytesVal)
			prevReq = append(prevReq, reqVal)
			prevNX = append(prevNX, nxVal)
			prevServfail = append(prevServfail, servfailVal)
			prevNoError = append(prevNoError, noErrorVal)

		}

		reqScore := zScore(currentReq, prevReq)
		bytesScore := zScore(currentBytes, prevBytes)
		nxScore := zScore(currentNX, prevNX)
		servfailScore := zScore(currentServfail, prevServfail)
		noErrorScore := zScore(currentNoError, prevNoError)

		finalScore, reason := maxScore(bytesScore, reqScore, nxScore, servfailScore, noErrorScore)

		/*
			NXDOMAIN_oranı = nx_domain_count / request_count	İLERİDE BUNA ÇEVİRİLEBİLİR
			SERVFAIL_oranı = servfail_count / request_count		EN MANTIKLISINI UYGULA
		*/

		isAnomaly, anomalyDegree := classifyAnomaly(finalScore, currentReq, currentBytes)

		anomalyCheck := &model.DomainAnomalyCheck{
			BucketStart:   bucketStart,
			Domain:        domain,
			Score:         finalScore,
			IsAnomaly:     isAnomaly,
			Severity:      string(anomalyDegree),
			Reason:        model.AnomalyReason(reason),
			TotalBytes:    int64(currentBytes),
			RequestCount:  int64(currentReq),
			BytesScore:    bytesScore,
			RequestScore:  reqScore,
			NXScore:       nxScore,
			ServfailScore: servfailScore,
		}

		if err = s.dc.Create(ctx, anomalyCheck); err != nil {
			return errors.New("anomalyService DomainCheck --> Create error: " + err.Error())
		}

		if !isAnomaly {
			continue
		}

		err = s.deepAnalyzeBucket(ctx, bucketStart, domain, reason, finalScore)
		if err != nil {
			return errors.New("anomalyService DeepAnalyzeBucket error: " + err.Error())
		}
	}
	return nil
}

func (s *anomalyService) deepAnalyzeBucket(ctx context.Context, bucketStart time.Time, domain, reason string, finalScore float64) error {
	buckets, err := s.br.BucketByDomainAndStartTime(ctx, domain, bucketStart)
	if err != nil {
		return errors.New("anomalyService BucketByDomainAndStartTime error: " + err.Error())
	}
	anomalyReason := model.AnomalyReason(reason)

	var bucket *model.TrafficBucket

	switch anomalyReason {
	case model.ReasonRequestSpike:
		bucket = forRequestSpike(buckets)
	case model.ReasonBytesSpike:
		bucket = forBytesSpike(buckets)
	case model.ReasonNXDomainSpike:
		bucket = forNXSpike(buckets)
	case model.ReasonServfailSpike:
		bucket = forServerfailSpike(buckets)
	case model.ReasonNoErrorSpike:
		bucket = forNoErrorSpike(buckets)
	}

	if bucket == nil {
		return nil
	}

	anomalyEvent := &model.AnomalyEvent{
		BucketStart: bucketStart,

		Domain:   domain,
		SourceIP: bucket.SourceIP,

		Score:     finalScore,
		IsAnomaly: true,
		Reason:    anomalyReason,

		TotalBytes:   bucket.TotalBytesSum,
		RequestCount: bucket.RequestCount,

		NXDomainCount: bucket.NXDomainCount,
		ServfailCount: bucket.ServfailCount,
		NoErrorCount:  bucket.NoErrorCount,

		// QueryType --> BURAYA ULAŞMIYOR --> bunu sor gerekiyorsa çöz
		Protocol: bucket.Protocol,

		Country: bucket.Country,
		ASN:     bucket.ASN,
	}

	err = s.ar.CreateAnomalyEvent(ctx, anomalyEvent)
	if err != nil {
		return errors.New("anomalyService CreateAnomalyEvent error: " + err.Error())
	}

	return nil
}

func (s *anomalyService) GetAnomalyEvents(ctx context.Context, vm *viewmodel.ExclusionRequest) (*viewmodel.ExclusionResponse, error) {
	domain := vm.Domain
	start, err := vm.StartTime()
	if err != nil {
		return nil, errors.New("anomalyService GetAnomalyEvents StartTime error: " + err.Error())
	}

	end, err := vm.EndTime()
	if err != nil {
		return nil, errors.New("anomalyService GetAnomalyEvents EndTime error: " + err.Error())
	}

	anomalyEvents, err := s.ar.GetEventsWithStartAndEndTime(ctx, domain, start, end)
	if err != nil {
		return nil, errors.New("anomalyService GetAnomalyEvents error: " + err.Error())
	}

	resp := help(anomalyEvents, domain)
	resp.StartDate = start.Format(time.RFC3339)
	resp.EndDate = end.Format(time.RFC3339)

	return resp, nil
}
