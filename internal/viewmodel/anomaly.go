package viewmodel

import "traffic-guarder/internal/model"

type AnomalyEvent struct {
	AttackStarted string `json:"attack_started"`
	AttackEnded   string `json:"attack_ended"`

	Domain string `json:"domain"`
	IP     string `json:"ip"`

	Score      float64 `json:"score"`
	MostReason string  `json:"most_reason"`

	TotalBytes    int64 `json:"total_bytes"`
	TotalRequests int64 `json:"total_requests"`

	TotalNXDomain int64 `json:"total_nxdomain"`
	TotalServfail int64 `json:"total_servfail"`
	TotalNoError  int64 `json:"total_no_error"`

	Protocol string `json:"protocol"`
	Country  string `json:"country"`
	ASN      string `json:"asn"`
}

func ToAnomalyEvent(event *model.AnomalyEvent) *AnomalyEvent {
	return &AnomalyEvent{}
}
