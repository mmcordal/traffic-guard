package main

import (
	"context"
	"log"
	"traffic-guarder/internal/simulator"
)

func main() {

	/*
		Mode'u seç uygun config Run da seçiliyor.
		"normal" --> simulator.NormalCfg
		"request_spike" --> simulator.RequestSpikeCfg
		"bytes_spike" --> simulator.BytesSpikeCfg
		"nx_domain_spike" --> simulator.NXDomainSpikeCfg
		"servfail_spike" --> simulator.ServfailSpikeCfg
	*/
	err := simulator.Run(context.Background(), "normal")
	if err != nil {
		log.Fatal(err)
	}
}
