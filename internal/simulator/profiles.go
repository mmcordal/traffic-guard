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

var protocol = []WeightedString{
	{Value: "UDP", Weight: 70},
	{Value: "TCP", Weight: 20},
	{Value: "DOH", Weight: 5},
	{Value: "DOT", Weight: 5},
}
