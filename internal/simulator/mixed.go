package simulator

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var mixedAttackModes = []string{
	"request_spike",
	"bytes_spike",
	"nxdomain_spike",
	"servfail_spike",
}

func RunMixed(ctx context.Context) error {
	fmt.Println("starting mixed simulator")

	mixedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)

	go func() {
		fmt.Println("normal traffic started")
		cfg := NormalCfg
		cfg.Duration = 0
		errCh <- runConfig(mixedCtx, cfg)
	}()

	go func() {
		errCh <- runAttackBursts(mixedCtx)
	}()

	err := <-errCh
	cancel()

	secondErr := <-errCh
	fmt.Println("mixed simulator stopped")

	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	if secondErr != nil && !errors.Is(secondErr, context.Canceled) {
		return secondErr
	}
	return nil
}

func runAttackBursts(ctx context.Context) error {
	for {
		waitDuration := randomDuration(60, 90)
		if err := waitForDuration(ctx, waitDuration); err != nil {
			return err
		}

		mode := mixedAttackModes[rand.Intn(len(mixedAttackModes))]
		fmt.Printf("selected attack mode: %s\n", mode)

		cfg, err := configForMode(mode)
		if err != nil {
			return err
		}
		cfg.Duration = randomDuration(30, 60)

		fmt.Printf("attack started: %s for %s\n", mode, cfg.Duration)
		if err := runConfig(ctx, cfg); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return fmt.Errorf("attack %s failed: %w", mode, err)
		}
		fmt.Printf("attack finished: %s\n", mode)
	}
}

func randomDuration(minSeconds, maxSeconds int64) time.Duration {
	return time.Duration(RandomInt64(minSeconds, maxSeconds)) * time.Second
}

func waitForDuration(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
