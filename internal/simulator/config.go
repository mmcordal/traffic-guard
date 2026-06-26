package simulator

import "time"

type Config struct {
	URL      string
	Domain   string
	Mode     string
	RPS      int
	Duration time.Duration
	Seed     int64
}

var NormalCfg = Config{ // rps * duration = 30 log
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "normal",
	RPS:      2,
	Duration: 15 * time.Second,
}

var RequestSpikeCfg = Config{ // rps * duration = 300 log
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr", // bi domaine yüklenmek istersen ilgili domaini gir.
	Mode:     "request_spike",
	RPS:      5,
	Duration: 60 * time.Second,
}

var BytesSpikeCfg = Config{ // rps * duration = 300 log
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "bytes_spike",
	RPS:      5,
	Duration: 60 * time.Second,
}

var NXDomainSpikeCfg = Config{ // rps * duration = 300 log
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "nx_domain_spike",
	RPS:      5,
	Duration: 60 * time.Second,
}

var ServfailSpikeCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "servfail_spike",
	RPS:      5,
	Duration: 60 * time.Second,
}
