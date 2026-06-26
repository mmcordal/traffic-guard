package cron

import (
	"context"
	"fmt"
	"log"
	"time"
	"traffic-guarder/internal/infrastructure/config"
	"traffic-guarder/internal/service"

	"github.com/robfig/cron/v3"
)

func Start(s service.AnomalyService, cfg config.AnalyzeConfig) {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	c := cron.New(cron.WithLocation(loc), cron.WithSeconds())

	schedule := fmt.Sprintf("@every %s", cfg.AnalyzeEvery())

	_, err := c.AddFunc(schedule, func() { // "*/1 * * * *" --> her dakikada || "0 17 * * *" --> her gün saat 17 de
		window := cfg.BucketWindow()
		bucketStart := time.Now().Truncate(window).Add(-window)

		log.Println("cron tetiklendi kral:", bucketStart)

		err := s.AnalyzeCompletedBucket(context.Background(), bucketStart)
		if err != nil {
			log.Println("cron error:", err)
		} else {
			log.Println("cron: AnalyzeCompletedBucket started!")
		}
	})

	if err != nil {
		log.Println("cron setup error:", err)
	}

	c.Start()
}
