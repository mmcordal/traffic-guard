package simulator

type TrafficLogPayload struct {
	Timestamp     string `json:"timestamp"`
	Domain        string `json:"domain"`
	QueryName     string `json:"query_name"`
	SourceIP      string `json:"source_ip"`
	QueryType     string `json:"query_type"`
	ResponseCode  string `json:"response_code"`
	RequestBytes  int64  `json:"request_bytes"`
	ResponseBytes int64  `json:"response_bytes"`
	TotalBytes    int64  `json:"total_bytes"`
	Protocol      string `json:"protocol"`
	LatencyMs     int64  `json:"latency_ms"`
	Country       string `json:"country"`
	ASN           string `json:"asn"`
}
