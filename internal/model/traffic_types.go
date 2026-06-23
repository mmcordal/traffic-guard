package model

import (
	"errors"
	"strings"
)

type QueryType string

const (
	QueryTypeA     QueryType = "A"
	QueryTypeAAAA  QueryType = "AAAA"
	QueryTypeMX    QueryType = "MX"
	QueryTypeTXT   QueryType = "TXT"
	QueryTypeCNAME QueryType = "CNAME"
	QueryTypeNS    QueryType = "NS"
	QueryTypeSOA   QueryType = "SOA"
	QueryTypePTR   QueryType = "PTR"
	QueryTypeANY   QueryType = "ANY"
	QueryTypeOther QueryType = "OTHER"
)

var validQueryTypes = map[QueryType]struct{}{
	QueryTypeA:     {},
	QueryTypeAAAA:  {},
	QueryTypeMX:    {},
	QueryTypeTXT:   {},
	QueryTypeCNAME: {},
	QueryTypeNS:    {},
	QueryTypeSOA:   {},
	QueryTypePTR:   {},
	QueryTypeANY:   {},
	QueryTypeOther: {},
}

func IsValidQueryType(value string) (QueryType, error) {
	normalized := QueryType(strings.ToUpper(strings.TrimSpace(value)))
	if _, ok := validQueryTypes[normalized]; !ok {
		return "", errors.New("invalid query type")
	}
	return normalized, nil
}

type ResponseCode string

const (
	ResponseCodeNoError  ResponseCode = "NOERROR"
	ResponseCodeNXDomain ResponseCode = "NXDOMAIN"
	ResponseCodeServfail ResponseCode = "SERVFAIL"
	ResponseCodeRefused  ResponseCode = "REFUSED"
	ResponseCodeFormerr  ResponseCode = "FORMERR"
)

var validResponseCodes = map[ResponseCode]struct{}{
	ResponseCodeNoError:  {},
	ResponseCodeNXDomain: {},
	ResponseCodeServfail: {},
	ResponseCodeRefused:  {},
	ResponseCodeFormerr:  {},
}

func IsValidResponseCode(value string) (ResponseCode, error) {
	normalized := ResponseCode(strings.ToUpper(strings.TrimSpace(value)))
	if _, ok := validResponseCodes[normalized]; !ok {
		return "", errors.New("invalid response code")
	}
	return normalized, nil
}

type Protocol string

const (
	ProtocolUDP Protocol = "UDP"
	ProtocolTCP Protocol = "TCP"
	ProtocolDoH Protocol = "DOH"
	ProtocolDoT Protocol = "DOT"
)

var validProtocols = map[Protocol]struct{}{
	ProtocolUDP: {},
	ProtocolTCP: {},
	ProtocolDoH: {},
	ProtocolDoT: {},
}

func IsValidProtocol(value string) (Protocol, error) {
	normalized := Protocol(strings.ToUpper(strings.TrimSpace(value)))
	if _, ok := validProtocols[normalized]; !ok {
		return "", errors.New("invalid protocol")
	}
	return normalized, nil
}

type AnomalyReason string

const (
	ReasonBytesSpike         AnomalyReason = "bytes_spike"
	ReasonRequestSpike       AnomalyReason = "request_spike"
	ReasonNXDomainSpike      AnomalyReason = "nxdomain_spike"
	ReasonServfailSpike      AnomalyReason = "servfail_spike"
	ReasonNoErrorSpike       AnomalyReason = "no_error_spike"
	ReasonHighIPContribution AnomalyReason = "high_ip_contribution"
	ReasonLowAndSlowPattern  AnomalyReason = "low_and_slow_pattern"
	ReasonFingerprintChange  AnomalyReason = "fingerprint_change"
)
