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
