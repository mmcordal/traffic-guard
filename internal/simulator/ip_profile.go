package simulator

import (
	"math/rand"
	"sync"
	"time"
)

type IPProfile struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	ASN         string `json:"asn"`
	Provider    string `json:"provider"`
	Kind        string `json:"kind"`
	Weight      int    `json:"weight"`
}

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu  sync.Mutex
)

var CommonDNSIPProfiles = []IPProfile{
	// Public DNS / Recursive resolvers
	{IP: "8.8.8.8", CountryCode: "US", ASN: "AS15169", Provider: "Google Public DNS", Kind: "public_resolver", Weight: 20},
	{IP: "8.8.4.4", CountryCode: "US", ASN: "AS15169", Provider: "Google Public DNS", Kind: "public_resolver", Weight: 14},

	{IP: "1.1.1.1", CountryCode: "US", ASN: "AS13335", Provider: "Cloudflare DNS", Kind: "public_resolver", Weight: 20},
	{IP: "1.0.0.1", CountryCode: "US", ASN: "AS13335", Provider: "Cloudflare DNS", Kind: "public_resolver", Weight: 14},

	{IP: "9.9.9.9", CountryCode: "US", ASN: "AS19281", Provider: "Quad9", Kind: "security_resolver", Weight: 10},
	{IP: "149.112.112.112", CountryCode: "US", ASN: "AS19281", Provider: "Quad9", Kind: "security_resolver", Weight: 7},

	{IP: "208.67.222.222", CountryCode: "US", ASN: "AS36692", Provider: "OpenDNS", Kind: "public_resolver", Weight: 9},
	{IP: "208.67.220.220", CountryCode: "US", ASN: "AS36692", Provider: "OpenDNS", Kind: "public_resolver", Weight: 7},

	{IP: "94.140.14.14", CountryCode: "CY", ASN: "AS212772", Provider: "AdGuard DNS", Kind: "filtering_resolver", Weight: 5},
	{IP: "94.140.15.15", CountryCode: "CY", ASN: "AS212772", Provider: "AdGuard DNS", Kind: "filtering_resolver", Weight: 5},

	{IP: "77.88.8.8", CountryCode: "RU", ASN: "AS13238", Provider: "Yandex DNS", Kind: "public_resolver", Weight: 4},
	{IP: "77.88.8.1", CountryCode: "RU", ASN: "AS13238", Provider: "Yandex DNS", Kind: "public_resolver", Weight: 3},

	// TR ISP benzeri simulator IP'leri
	// Bunlar gerçek DNS resolver IP doğrulaması için değil, trafik dağılımı simülasyonu için.
	{IP: "85.100.10.15", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 18},
	{IP: "85.105.25.41", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 16},

	{IP: "176.232.44.12", CountryCode: "TR", ASN: "AS16135", Provider: "Turkcell", Kind: "mobile_isp", Weight: 12},
	{IP: "176.233.92.80", CountryCode: "TR", ASN: "AS16135", Provider: "Turkcell", Kind: "mobile_isp", Weight: 10},

	{IP: "212.156.80.21", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 10},
	{IP: "212.156.92.33", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 8},

	{IP: "78.186.72.9", CountryCode: "TR", ASN: "AS15897", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 8},
	{IP: "78.187.60.17", CountryCode: "TR", ASN: "AS15897", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 7},

	{IP: "94.54.20.18", CountryCode: "TR", ASN: "AS47524", Provider: "Turksat Kablo", Kind: "isp_resolver", Weight: 6},
	{IP: "94.55.70.26", CountryCode: "TR", ASN: "AS47524", Provider: "Turksat Kablo", Kind: "isp_resolver", Weight: 5},

	{IP: "31.223.12.77", CountryCode: "TR", ASN: "AS12735", Provider: "TurkNet", Kind: "isp_resolver", Weight: 6},
	{IP: "31.223.45.91", CountryCode: "TR", ASN: "AS12735", Provider: "TurkNet", Kind: "isp_resolver", Weight: 5},

	// Şüpheli / saldırı simülasyonu için hosting benzeri profiller
	{IP: "45.155.205.12", CountryCode: "NL", ASN: "AS9009", Provider: "M247", Kind: "hosting_suspicious", Weight: 2},
	{IP: "185.220.101.44", CountryCode: "DE", ASN: "AS60729", Provider: "Tor/Hosting Exit", Kind: "hosting_suspicious", Weight: 2},
	{IP: "198.98.51.189", CountryCode: "US", ASN: "AS53667", Provider: "FranTech", Kind: "hosting_suspicious", Weight: 2},
	{IP: "91.200.12.66", CountryCode: "UA", ASN: "AS48666", Provider: "Suspicious Network", Kind: "hosting_suspicious", Weight: 1},
}

func PickIPProfile() IPProfile {
	mu.Lock()
	defer mu.Unlock()

	totalWeight := 0

	for _, profile := range CommonDNSIPProfiles {
		totalWeight += profile.Weight
	}

	if totalWeight <= 0 {
		return CommonDNSIPProfiles[0]
	}

	randomWeight := rng.Intn(totalWeight)

	for _, profile := range CommonDNSIPProfiles {
		randomWeight -= profile.Weight
		if randomWeight < 0 {
			return profile
		}
	}

	return CommonDNSIPProfiles[len(CommonDNSIPProfiles)-1]
}

func PickSuspiciousIPProfile() IPProfile {
	suspicious := make([]IPProfile, 0)

	for _, profile := range CommonDNSIPProfiles {
		if profile.Kind == "hosting_suspicious" {
			suspicious = append(suspicious, profile)
		}
	}

	if len(suspicious) == 0 {
		return PickIPProfile()
	}

	mu.Lock()
	defer mu.Unlock()

	return suspicious[rng.Intn(len(suspicious))]
}
