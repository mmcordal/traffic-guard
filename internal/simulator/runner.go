package simulator

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Run(ctx context.Context, mode string) error {
	if mode == "mixed" {
		return RunMixed(ctx)
	}

	cfg, err := configForMode(mode)
	if err != nil {
		return err
	}

	if cfg.Duration <= 0 {
		cfg.Duration = time.Second * 10
	}

	return runConfig(ctx, cfg)
}

func configForMode(mode string) (Config, error) {
	switch mode {
	case "normal":
		return NormalCfg, nil
	case "request_spike":
		return RequestSpikeCfg, nil
	case "bytes_spike":
		return BytesSpikeCfg, nil
	case "nxdomain_spike", "nx_domain_spike":
		cfg := NXDomainSpikeCfg
		cfg.Mode = "nxdomain_spike"
		return cfg, nil
	case "servfail_spike":
		return ServfailSpikeCfg, nil
	default:
		return Config{}, fmt.Errorf("unknown simulator mode: %s", mode)
	}
}

func runConfig(ctx context.Context, cfg Config) error {
	if cfg.URL == "" {
		return errors.New("simulator URL is required")
	}

	if cfg.RPS <= 0 {
		cfg.RPS = 1
	}

	interval := time.Second / time.Duration(cfg.RPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	runCtx := ctx
	cancel := func() {}
	if cfg.Duration > 0 {
		runCtx, cancel = context.WithTimeout(ctx, cfg.Duration)
	}
	defer cancel()

	sendCount := 0

	for {
		select {
		case <-runCtx.Done():
			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return ctx.Err()
			}
			fmt.Printf("simulator finished, sent count: %d\n", sendCount)
			return nil
		case <-ticker.C:
			log := generateLog(cfg.Mode, cfg.Domain)

			err := SendLog(runCtx, cfg.URL, log)
			if err != nil {
				if runCtx.Err() != nil {
					return runCtx.Err()
				}
				return fmt.Errorf("number of logs sent so far: %d simulator failed to send log: %w", sendCount, err)
			}
			sendCount++
		}
	}
}

func generateLog(mode, domain string) TrafficLogPayload {
	switch mode {
	case "normal":
		return GenerateNormalLog(domain)
	case "request_spike":
		return GenerateRequestSpikeLog(domain)
	case "bytes_spike":
		return GenerateBytesSpikeLog(domain)
	case "nxdomain_spike", "nx_domain_spike":
		return GenerateNXDomainSpikeLog(domain)
	case "servfail_spike":
		return GenerateServfailSpikeLog(domain)
	default:
		return GenerateNormalLog(domain)
	}
}
