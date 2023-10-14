package main

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

type config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "URI:port")
	flag.IntVar(&flagReportInterval, "r", 10, "interval to report")
	flag.IntVar(&flagPollInterval, "p", 2, "interval to poll")
	flag.Parse()
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Address != "" {
		flagRunAddr = cfg.Address
	}
	if cfg.ReportInterval != 0 {
		flagReportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		flagPollInterval = cfg.PollInterval
	}
}
