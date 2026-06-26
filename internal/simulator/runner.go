package simulator

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Run(ctx context.Context, mode string) error {
	cfg := Config{}

	switch mode {
	case "normal":
		cfg = NormalCfg
	case "request_spike":
		cfg = RequestSpikeCfg
	case "bytes_spike":
		cfg = BytesSpikeCfg
	case "nx_domain_spike":
		cfg = NXDomainSpikeCfg
	case "servfail_spike":
		cfg = ServfailSpikeCfg
	}

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
			log := TrafficLogPayload{}

			switch cfg.Mode {
			case "normal":
				log = GenerateNormalLog(cfg.Domain)
			case "request_spike":
				log = GenerateRequestSpikeLog(cfg.Domain)
			case "bytes_spike":
				log = GenerateBytesSpikeLog(cfg.Domain)
			case "nx_domain_spike":
				log = GenerateNXDomainSpikeLog(cfg.Domain)
			case "servfail_spike":
				log = GenerateServfailSpikeLog(cfg.Domain)
			default:
				log = GenerateNormalLog(cfg.Domain)
			}

			err := SendLog(timeoutCtx, cfg.URL, log)
			if err != nil {
				return fmt.Errorf("number of logs sent so far: %d simulator failed to send log: %w", sendCount, err)
			}
			sendCount++
		}
	}

}
