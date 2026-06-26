package simulator

import "time"

func GenerateRequestSpikeLog(domain string) TrafficLogPayload {
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

func GenerateBytesSpikeLog(domain string) TrafficLogPayload {
	if domain == "" {
		domain = Pick(sampleDomains)
	}

	requestBytes := RandomInt64(60, 150)
	responseBytes := RandomInt64(1500, 10000)

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

func GenerateNXDomainSpikeLog(domain string) TrafficLogPayload {
	if domain == "" {
		domain = Pick(sampleDomains)
	}

	requestBytes := RandomInt64(60, 150)
	responseBytes := RandomInt64(200, 900)

	ipProfile := WeightedPickIP(CommonDNSIPProfiles)

	return TrafficLogPayload{
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Domain:        domain,
		QueryName:     RandomAlphaNumeric(int(RandomInt64(8, 12))) + "." + domain,
		SourceIP:      ipProfile.IP,
		QueryType:     WeightedPickString(queryTypeforNXSpike),
		ResponseCode:  "NXDOMAIN",
		RequestBytes:  requestBytes,
		ResponseBytes: responseBytes,
		TotalBytes:    requestBytes + responseBytes,
		Protocol:      WeightedPickString(protocols),
		LatencyMs:     RandomInt64(100, 500),
		Country:       ipProfile.CountryCode,
		ASN:           ipProfile.ASN,
	}
}

func GenerateServfailSpikeLog(domain string) TrafficLogPayload {
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
		ResponseCode:  "SERVFAIL",
		RequestBytes:  requestBytes,
		ResponseBytes: responseBytes,
		TotalBytes:    requestBytes + responseBytes,
		Protocol:      WeightedPickString(protocols),
		LatencyMs:     RandomInt64(2000, 6000),
		Country:       ipProfile.CountryCode,
		ASN:           ipProfile.ASN,
	}
}
