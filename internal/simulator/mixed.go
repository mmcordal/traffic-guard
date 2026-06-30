package simulator

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	mixedInitialNormalDuration = 10 * time.Minute
	mixedWaveCooldownDuration  = 10 * time.Minute
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
	fmt.Printf("initial normal-only period started for %s\n", mixedInitialNormalDuration)
	if err := waitForDuration(ctx, mixedInitialNormalDuration); err != nil {
		return err
	}

	for {
		if err := runAttackWave(ctx); err != nil {
			return err
		}

		fmt.Printf("attack wave finished; waiting %s before next wave\n", mixedWaveCooldownDuration)
		if err := waitForDuration(ctx, mixedWaveCooldownDuration); err != nil {
			return err
		}
	}
}

func runAttackWave(ctx context.Context) error {
	attackCount := int(RandomInt64(2, 3))
	fmt.Printf("attack wave started with %d attacks\n", attackCount)

	for attackIndex := 0; attackIndex < attackCount; attackIndex++ {
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

		if attackIndex == attackCount-1 {
			continue
		}

		gapDuration := randomDuration(10, 20)
		fmt.Printf("waiting %s before next attack in wave\n", gapDuration)
		if err := waitForDuration(ctx, gapDuration); err != nil {
			return err
		}
	}

	return nil
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
