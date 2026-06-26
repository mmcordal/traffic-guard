package simulator

import (
	"time"
)

func GenerateNormalLog(domain string) TrafficLogPayload {
	if domain == "" {
		domain = Pick(sampleDomains)
	}

	requestBytes := RandomInt64(60, 150)
	responseBytes := RandomInt64(200, 900)

	ipProfile := WeightedPickIP(CommonDNSIPProfiles)

	return TrafficLogPayload{
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Domain:        domain,
		QueryName:     RandomReadableWord() + "." + domain,
		SourceIP:      ipProfile.IP,
		QueryType:     WeightedPickString(queryType),
		ResponseCode:  WeightedPickString(responseCodes),
		RequestBytes:  requestBytes,
		ResponseBytes: responseBytes,
		TotalBytes:    requestBytes + responseBytes,
		Protocol:      WeightedPickString(protocols),
		LatencyMs:     RandomInt64(10, 120),
		Country:       ipProfile.CountryCode,
		ASN:           ipProfile.ASN,
	}

}
