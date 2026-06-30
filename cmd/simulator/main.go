package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"traffic-guarder/internal/simulator"
)

func main() {
	mode := "normal"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := simulator.Run(ctx, mode); err != nil {
		log.Fatal(err)
	}
}
