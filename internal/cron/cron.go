package cron

import (
	"context"
	"log"
	"time"
	"traffic-guarder/internal/service"

	"github.com/robfig/cron/v3"
)

func Start(s service.AnomalyService) {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("*/1 * * * *", func() { // "*/1 * * * *" --> her dakikada || "0 17 * * *" --> her gün saat 17 de
		log.Println("cron tetiklendi kral")

		err := s.AnalyzeCompletedBucket(context.Background(), time.Now().Truncate(time.Minute).Add(-1*time.Minute))

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
