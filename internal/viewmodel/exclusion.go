package viewmodel

import (
	"errors"
	"strings"
	"time"
)

type ExclusionRequest struct {
	StartDate string `json:"startDate" validate:"required,datetime=2006-01-02T15:04:05Z07:00" labelName:"startDate"`
	EndDate   string `json:"endDate" validate:"required,datetime=2006-01-02T15:04:05Z07:00" labelName:"endDate"`
	Domain    string `json:"domain" validate:"required,max=255" labelName:"domain"`
}

type ExclusionResponse struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	Domain string `json:"domain"`

	AnomalyEvents []*AnomalyEvent `json:"anomaly_events"`
	DangerousIPs  []string        `json:"dangerous_ips"`

	SpentTotalBytes int64 `json:"spent_total_bytes"`
	TotalRequests   int64 `json:"total_requests"`

	TotalNXDomain int64 `json:"total_nxdomain"`
	TotalServfail int64 `json:"total_servfail"`
	TotalNoError  int64 `json:"total_no_error"`

	AverageScore float64 `json:"average_score"`

	MostAnomalyReason string `json:"most_anomaly_reason"`
	MostProtocol      string `json:"most_protocol"`
	MostCountry       string `json:"most_country"`
	MostASN           string `json:"most_asn"`
}

func (r ExclusionRequest) Validate() []error {
	var errs []error

	start, startErr := time.Parse(time.RFC3339, r.StartDate)
	if startErr != nil {
		errs = append(errs, errors.New("Başlangıç tarihi uygun bir tarih formatında olmalı"))
	}

	end, endErr := time.Parse(time.RFC3339, r.EndDate)
	if endErr != nil {
		errs = append(errs, errors.New("Bitiş tarihi uygun bir tarih formatında olmalı"))
	}

	if startErr == nil && endErr == nil {
		if start.After(end) {
			errs = append(errs, errors.New("Başlangıç tarihi bitiş tarihinden sonra olamaz"))
		}
	}

	if !isValidDomain(r.Domain) {
		errs = append(errs, errors.New("Domain geçerli bir domain olmalı!"))
	}

	return errs
}

func isValidDomain(domain string) bool {
	domain = strings.TrimSpace(domain)

	if domain == "" {
		return false
	}

	if len(domain) > 255 {
		return false
	}

	if strings.Contains(domain, "://") {
		return false
	}

	if strings.Contains(domain, " ") {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

func (r ExclusionRequest) StartTime() (time.Time, error) {
	return time.Parse(time.RFC3339, r.StartDate)
}

func (r ExclusionRequest) EndTime() (time.Time, error) {
	return time.Parse(time.RFC3339, r.EndDate)
}
