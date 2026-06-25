package main

import (
	"context"
	"log"
	"time"
	"traffic-guarder/internal/simulator"
)

func main() {

	cfg := simulator.Config{
		URL:      "http://localhost:8080/api/v1/traffic-log/",
		Domain:   "example.com",
		Mode:     "normal",
		RPS:      2,
		Duration: 15 * time.Second,
	}
	err := simulator.Run(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
}
