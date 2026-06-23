package viewmodel

import (
	"errors"
	"time"
	"traffic-guarder/internal/model"
)

type CreateTrafficLogRequest struct {
	Timestamp string `json:"timestamp" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" labelName:"zaman"`

	Domain    string `json:"domain" validate:"required,max=255" labelName:"domain"`
	QueryName string `json:"query_name" validate:"required,max=255" labelName:"sorgu adı"`
	SourceIP  string `json:"source_ip" validate:"required,ip,max=45" labelName:"kaynak ip"`

	QueryType    string `json:"query_type" validate:"required" labelName:"sorgu tipi"`
	ResponseCode string `json:"response_code" validate:"required" labelName:"cevap kodu"`

	RequestBytes  int64 `json:"request_bytes" validate:"gte=0" labelName:"istek byte"`
	ResponseBytes int64 `json:"response_bytes" validate:"gte=0" labelName:"cevap byte"`
	TotalBytes    int64 `json:"total_bytes" validate:"gte=0" labelName:"toplam byte"`

	Protocol  string `json:"protocol" validate:"required" labelName:"protokol"`
	LatencyMs int64  `json:"latency_ms" validate:"gte=0" labelName:"gecikme"`

	Country string `json:"country" validate:"omitempty,len=2" labelName:"ülke"`
	ASN     string `json:"asn" validate:"omitempty,max=32" labelName:"asn"`
}

func (r CreateTrafficLogRequest) Validate() []error {
	var errs []error

	if _, err := model.IsValidQueryType(r.QueryType); err != nil {
		errs = append(errs, errors.New("geçersiz query_type"))
	}

	if _, err := model.IsValidResponseCode(r.ResponseCode); err != nil {
		errs = append(errs, errors.New("geçersiz response_code"))
	}

	if _, err := model.IsValidProtocol(r.Protocol); err != nil {
		errs = append(errs, errors.New("geçersiz protocol"))
	}

	if r.TotalBytes == 0 && r.RequestBytes+r.ResponseBytes == 0 {
		errs = append(errs, errors.New("total_bytes veya request_bytes + response_bytes değerlerinden biri dolu olmalı"))
	}

	return errs
}

func (r CreateTrafficLogRequest) ToModel() (*model.TrafficLog, error) {
	queryType, err := model.IsValidQueryType(r.QueryType)
	if err != nil {
		return nil, err
	}

	responseCode, err := model.IsValidResponseCode(r.ResponseCode)
	if err != nil {
		return nil, err
	}

	protocol, err := model.IsValidProtocol(r.Protocol)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()
	if r.Timestamp != "" {
		parsedTime, err := time.Parse(time.RFC3339, r.Timestamp)
		if err != nil {
			return nil, err
		}
		timestamp = parsedTime
	}

	totalBytes := r.TotalBytes
	if totalBytes == 0 {
		totalBytes = r.RequestBytes + r.ResponseBytes
	}

	return &model.TrafficLog{
		Timestamp:     timestamp,
		Domain:        r.Domain,
		QueryName:     r.QueryName,
		SourceIP:      r.SourceIP,
		QueryType:     queryType,
		ResponseCode:  responseCode,
		RequestBytes:  r.RequestBytes,
		ResponseBytes: r.ResponseBytes,
		TotalBytes:    totalBytes,
		Protocol:      protocol,
		LatencyMs:     r.LatencyMs,
		Country:       r.Country,
		ASN:           r.ASN,
	}, nil
}

type TrafficLogResponse struct {
	ID int64 `json:"id"`

	CreatedAt string `json:"created_at"`
	Timestamp string `json:"timestamp"`

	Domain    string `json:"domain"`
	QueryName string `json:"query_name"`
	SourceIP  string `json:"source_ip"`

	QueryType    model.QueryType    `json:"query_type"`
	ResponseCode model.ResponseCode `json:"response_code"`

	RequestBytes  int64 `json:"request_bytes"`
	ResponseBytes int64 `json:"response_bytes"`
	TotalBytes    int64 `json:"total_bytes"`

	Protocol  model.Protocol `json:"protocol"`
	LatencyMs int64          `json:"latency_ms"`

	Country string `json:"country,omitempty"`
	ASN     string `json:"asn,omitempty"`
}

func ToTrafficLogResponse(t *model.TrafficLog) *TrafficLogResponse {
	if t == nil {
		return nil
	}

	return &TrafficLogResponse{
		ID:        t.ID,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		Timestamp: t.Timestamp.Format(time.RFC3339),

		Domain:    t.Domain,
		QueryName: t.QueryName,
		SourceIP:  t.SourceIP,

		QueryType:    t.QueryType,
		ResponseCode: t.ResponseCode,

		RequestBytes:  t.RequestBytes,
		ResponseBytes: t.ResponseBytes,
		TotalBytes:    t.TotalBytes,

		Protocol:  t.Protocol,
		LatencyMs: t.LatencyMs,

		Country: t.Country,
		ASN:     t.ASN,
	}
}

func ToTrafficLogResponses(logs []*model.TrafficLog) []*TrafficLogResponse {
	out := make([]*TrafficLogResponse, 0, len(logs))

	for _, log := range logs {
		out = append(out, ToTrafficLogResponse(log))
	}

	return out
}
