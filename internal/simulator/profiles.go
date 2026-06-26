package simulator

type DomainProfile struct {
	Domain string
	Weight int
}

var domainProfiles = []DomainProfile{}

// 20 tane
var sampleDomains = []string{
	"example.com",
	"api.example.com",
	"cdn.example.com",
	"google.com",
	"youtube.com",
	"music.youtube.com",
	"github.com",
	"openai.com",
	"cloudflare.com",
	"trendyol.com",
	"hepsiburada.com",
	"pau.edu.tr",
	"rule-dns.com",
	"netinternet.tr",
	"medium.com",
	"go.dev",
	"bun.uptrace.dev",
	"redis.io",
	"linkedin.com",
	"ajet.com",
}

// 26 tane
var subDomains = []string{
	"www",
	"mail",
	"webmail",
	"blog",
	"shop",
	"store",
	"api",
	"admin",
	"cp",
	"cpanel",
	"dev",
	"staging",
	"test",
	"ftp",
	"ns1",
	"ns2",
	"m",
	"cdn",
	"auth",
	"static",
	"img",
	"app",
	"login",
	"logout",
	"payment",
	"dashboard",
}

type WeightedString struct {
	Value  string
	Weight int
}

var responseCodes = []WeightedString{
	{Value: "NOERROR", Weight: 85},
	{Value: "NXDOMAIN", Weight: 11},
	{Value: "SERVFAIL", Weight: 3},
	{Value: "REFUSED", Weight: 1},
}

var queryType = []WeightedString{
	{Value: "A", Weight: 70},
	{Value: "AAAA", Weight: 14},
	{Value: "MX", Weight: 2},
	{Value: "TXT", Weight: 3},
	{Value: "CNAME", Weight: 5},
	{Value: "NS", Weight: 2},
	{Value: "SOA", Weight: 1},
	{Value: "PTR", Weight: 1},
	{Value: "ANY", Weight: 1},
	{Value: "OTHER", Weight: 1},
}

var queryTypeforNXSpike = []WeightedString{
	{Value: "A", Weight: 80},
	{Value: "AAAA", Weight: 14},
	{Value: "MX", Weight: 2},
	{Value: "ANY", Weight: 2},
	{Value: "OTHER", Weight: 2},
}

var queryTypeforServfail = []WeightedString{
	{Value: "A", Weight: 80},
	{Value: "AAAA", Weight: 14},
	{Value: "MX", Weight: 2},
	{Value: "TXT", Weight: 2},
	{Value: "NS", Weight: 2},
}

var protocols = []WeightedString{
	{Value: "UDP", Weight: 70},
	{Value: "TCP", Weight: 20},
	{Value: "DOH", Weight: 5},
	{Value: "DOT", Weight: 5},
}

type IPProfile struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	ASN         string `json:"asn"`
	Provider    string `json:"provider"`
	Kind        string `json:"kind"`
	Weight      int    `json:"weight"`
}

var CommonDNSIPProfiles = []IPProfile{
	// Public DNS / Recursive resolvers
	{IP: "8.8.8.8", CountryCode: "US", ASN: "AS15169", Provider: "Google Public DNS", Kind: "public_resolver", Weight: 8},
	{IP: "8.8.4.4", CountryCode: "US", ASN: "AS15169", Provider: "Google Public DNS", Kind: "public_resolver", Weight: 6},

	{IP: "1.1.1.1", CountryCode: "US", ASN: "AS13335", Provider: "Cloudflare DNS", Kind: "public_resolver", Weight: 8},
	{IP: "1.0.0.1", CountryCode: "US", ASN: "AS13335", Provider: "Cloudflare DNS", Kind: "public_resolver", Weight: 6},

	{IP: "9.9.9.9", CountryCode: "US", ASN: "AS19281", Provider: "Quad9", Kind: "security_resolver", Weight: 5},
	{IP: "149.112.112.112", CountryCode: "US", ASN: "AS19281", Provider: "Quad9", Kind: "security_resolver", Weight: 4},

	{IP: "208.67.222.222", CountryCode: "US", ASN: "AS36692", Provider: "OpenDNS", Kind: "public_resolver", Weight: 5},
	{IP: "208.67.220.220", CountryCode: "US", ASN: "AS36692", Provider: "OpenDNS", Kind: "public_resolver", Weight: 4},

	{IP: "94.140.14.14", CountryCode: "CY", ASN: "AS212772", Provider: "AdGuard DNS", Kind: "filtering_resolver", Weight: 3},
	{IP: "94.140.15.15", CountryCode: "CY", ASN: "AS212772", Provider: "AdGuard DNS", Kind: "filtering_resolver", Weight: 3},

	{IP: "77.88.8.8", CountryCode: "RU", ASN: "AS13238", Provider: "Yandex DNS", Kind: "public_resolver", Weight: 2},
	{IP: "77.88.8.1", CountryCode: "RU", ASN: "AS13238", Provider: "Yandex DNS", Kind: "public_resolver", Weight: 2},

	// TR ISP benzeri simulator IP'leri
	// Normal trafik TR ağırlıklı gelsin diye weight'ler yüksek.
	{IP: "88.255.1.10", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 35},
	{IP: "88.255.1.11", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 32},
	{IP: "88.255.1.12", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 30},

	{IP: "85.100.10.15", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 35},
	{IP: "85.105.25.41", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 30},
	{IP: "176.88.14.77", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 24},
	{IP: "176.88.14.78", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 22},
	{IP: "176.88.14.79", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 20},
	{IP: "94.54.22.100", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 18},
	{IP: "94.54.22.101", CountryCode: "TR", ASN: "AS9121", Provider: "Turk Telekom", Kind: "isp_resolver", Weight: 16},

	{IP: "95.70.12.45", CountryCode: "TR", ASN: "AS47331", Provider: "Turk Telekom / TTNET", Kind: "isp_resolver", Weight: 28},
	{IP: "95.70.12.46", CountryCode: "TR", ASN: "AS47331", Provider: "Turk Telekom / TTNET", Kind: "isp_resolver", Weight: 24},
	{IP: "85.105.201.33", CountryCode: "TR", ASN: "AS47331", Provider: "Turk Telekom / TTNET", Kind: "isp_resolver", Weight: 20},

	{IP: "78.189.33.21", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 24},
	{IP: "78.189.33.22", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 22},
	{IP: "31.223.88.12", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 18},
	{IP: "37.155.44.10", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "mobile_isp", Weight: 16},
	{IP: "37.155.44.11", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "mobile_isp", Weight: 14},
	{IP: "212.156.80.21", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 22},
	{IP: "212.156.92.33", CountryCode: "TR", ASN: "AS34984", Provider: "Turkcell Superonline", Kind: "isp_resolver", Weight: 18},

	{IP: "46.2.145.90", CountryCode: "TR", ASN: "AS15924", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 22},
	{IP: "46.2.145.91", CountryCode: "TR", ASN: "AS15924", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 20},
	{IP: "212.174.55.19", CountryCode: "TR", ASN: "AS15924", Provider: "Vodafone Turkey", Kind: "isp_resolver", Weight: 16},
	{IP: "213.74.120.50", CountryCode: "TR", ASN: "AS15924", Provider: "Vodafone Turkey", Kind: "isp_resolver", Weight: 12},
	{IP: "213.74.120.51", CountryCode: "TR", ASN: "AS15924", Provider: "Vodafone Turkey", Kind: "isp_resolver", Weight: 10},

	{IP: "176.232.44.12", CountryCode: "TR", ASN: "AS16135", Provider: "Turkcell", Kind: "mobile_isp", Weight: 25},
	{IP: "176.233.92.80", CountryCode: "TR", ASN: "AS16135", Provider: "Turkcell", Kind: "mobile_isp", Weight: 22},

	{IP: "78.186.72.9", CountryCode: "TR", ASN: "AS15897", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 18},
	{IP: "78.187.60.17", CountryCode: "TR", ASN: "AS15897", Provider: "Vodafone Turkey", Kind: "mobile_isp", Weight: 15},

	{IP: "94.54.20.18", CountryCode: "TR", ASN: "AS47524", Provider: "Turksat Kablo", Kind: "isp_resolver", Weight: 16},
	{IP: "94.55.70.26", CountryCode: "TR", ASN: "AS47524", Provider: "Turksat Kablo", Kind: "isp_resolver", Weight: 14},

	{IP: "31.223.12.77", CountryCode: "TR", ASN: "AS12735", Provider: "TurkNet", Kind: "isp_resolver", Weight: 14},
	{IP: "31.223.45.91", CountryCode: "TR", ASN: "AS12735", Provider: "TurkNet", Kind: "isp_resolver", Weight: 12},

	{IP: "185.86.12.20", CountryCode: "TR", ASN: "AS42926", Provider: "Netinternet", Kind: "hosting_tr", Weight: 10},
	{IP: "185.86.12.21", CountryCode: "TR", ASN: "AS42926", Provider: "Netinternet", Kind: "hosting_tr", Weight: 8},

	{IP: "193.192.98.30", CountryCode: "TR", ASN: "AS8517", Provider: "Academic / University Network", Kind: "edu_network", Weight: 8},
	{IP: "193.192.98.31", CountryCode: "TR", ASN: "AS8517", Provider: "Academic / University Network", Kind: "edu_network", Weight: 6},

	// Şüpheli / saldırı simülasyonu için hosting benzeri profiller
	// Normal trafikte nadir gelsin diye weight düşük.
	{IP: "45.155.205.12", CountryCode: "NL", ASN: "AS9009", Provider: "M247", Kind: "hosting_suspicious", Weight: 1},
	{IP: "185.220.101.44", CountryCode: "DE", ASN: "AS60729", Provider: "Tor/Hosting Exit", Kind: "hosting_suspicious", Weight: 1},
	{IP: "198.98.51.189", CountryCode: "US", ASN: "AS53667", Provider: "FranTech", Kind: "hosting_suspicious", Weight: 1},
	{IP: "91.200.12.66", CountryCode: "UA", ASN: "AS48666", Provider: "Suspicious Network", Kind: "hosting_suspicious", Weight: 1},
}
