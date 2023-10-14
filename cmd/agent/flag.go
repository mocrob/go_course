package main

import "flag"

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "URI:port")
	flag.IntVar(&flagReportInterval, "r", 10, "interval to report")
	flag.IntVar(&flagPollInterval, "p", 2, "interval to poll")
	flag.Parse()
}
