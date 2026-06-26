package simulator

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Run(ctx context.Context, cfg Config) error {
	if cfg.URL == "" {
		return errors.New("simulator URL is required")
	}

	if cfg.RPS <= 0 {
		cfg.RPS = 1
	}

	if cfg.Duration <= 0 {
		cfg.Duration = time.Second * 10
	}

	interval := time.Second / time.Duration(cfg.RPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.Duration)
	defer cancel()

	sendCount := 0

	for {
		select {
		case <-timeoutCtx.Done():
			fmt.Printf("simulator finished, sent count: %d\n", sendCount)
			return nil
		case <-ticker.C:
			log := GenerateNormalLog(cfg.Domain)

			err := SendLog(timeoutCtx, cfg.URL, log)
			if err != nil {
				return fmt.Errorf("number of logs sent so far: %d simulator failed to send log: %w", sendCount, err)
			}
			sendCount++
		}
	}

}
