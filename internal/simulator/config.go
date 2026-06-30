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

var NormalCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "normal",
	RPS:      60,
	Duration: 330 * time.Second,
} // 2 * 120 = 240 log

var RequestSpikeCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr", // bi domaine yüklenmek istersen ilgili domaini gir.
	Mode:     "request_spike",
	RPS:      500,
	Duration: 30 * time.Second,
} // 10 * 60 = 600 log

// boş olacaksa da şöyle: // 20 * 60 = 1200 log
/*
	Domain:   "",
	RPS:      20,
	Duration: 60 * time.Second,
*/

var BytesSpikeCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "bytes_spike",
	RPS:      100,
	Duration: 60 * time.Second,
} // 5 * 60 = 300 log	-- 1500-10000 response bytes

// boş olacaksa da şöyle: // 10 * 60 = 600 log
/*
	Domain:   "",
	RPS:      10,
	Duration: 60 * time.Second,
*/

var NXDomainSpikeCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "nx_domain_spike",
	RPS:      100,
	Duration: 60 * time.Second,
} // 8 * 60 = 480 log

// boş olacaksa da şöyle: // 15 * 60 = 900 log
/*
	Domain:   "",
	RPS:      15,
	Duration: 60 * time.Second,
*/

var ServfailSpikeCfg = Config{
	URL:      "http://localhost:8080/api/v1/traffic-log/",
	Domain:   "netinternet.tr",
	Mode:     "servfail_spike",
	RPS:      100,
	Duration: 60 * time.Second,
}

// boş olacaksa da şöyle: // 10 * 60 = 600 log
/*
	Domain:   "",
	RPS:      10,
	Duration: 60 * time.Second,
*/
