package service

import (
	"context"
	"errors"
	"strconv"
	"time"
	"traffic-guarder/internal/infrastructure/cache"
	"traffic-guarder/internal/model"
	"traffic-guarder/internal/repository"
)

type AnomalyService interface {
	AnalyzeCompletedBucket(ctx context.Context, bucketStart time.Time) error
}

type anomalyService struct {
	br repository.BucketRepository
	dc repository.DomainCheck
	bc cache.BucketCache
}

func NewAnomalyService(dc repository.DomainCheck, bc cache.BucketCache, br repository.BucketRepository) AnomalyService {
	return &anomalyService{dc: dc, bc: bc, br: br}
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

		previousMinutes, err := s.bc.GetPreviousBucketMinutes(ctx, domain, bucketStart.Add(-1*time.Minute), 10)
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
			NXDOMAIN oranı = nx_domain_count / request_count	İLERİDE BUNA EVRİLEBİLİR
			SERVFAIL oranı = servfail_count / request_count		EN MANTIKLISINI UYGULA
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

		err = s.DeepAnalyzeBucket(ctx, bucketStart, domain, reason, finalScore)
		if err != nil {
			return errors.New("anomalyService DeepAnalyzeBucket error: " + err.Error())
		}
	}
	return nil
}

func (s *anomalyService) DeepAnalyzeBucket(ctx context.Context, bucketStart time.Time, domain, reason string, finalScore float64) error {
	var dangerousIPs []string

	buckets, err := s.br.BucketByDomainAndStart(ctx, domain, bucketStart)
	if err != nil {
		return errors.New("anomalyService BucketByDomainAndStart error: " + err.Error())
	}
	anomalyReason := model.AnomalyReason(reason)

	_ = buckets
	_ = dangerousIPs
	_ = anomalyReason

	/*

		var bucket model.TrafficBucket

		switch anomalyReason {
		case model.ReasonRequestSpike:
			bucket := forRequestSpike(buckets)
		case model.ReasonBytesSpike:
			bucket := forBytesSpike(buckets)
		case model.ReasonNXDomainSpike:
			bucket := forNXSpike(buckets)
		case model.ReasonServfailSpike:
			bucket := forServerfailSpike(buckets)
		case model.ReasonNoErrorSpike:
			bucket := forNoErrorSpike(buckets)
		}

		if anomalyReason == model.ReasonRequestSpike {
			bucket := forRequestSpike(buckets)
		}

		anomalyEvent := &model.AnomalyEvent{}

	*/
	return nil
}
