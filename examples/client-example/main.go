package main

import (
	"fmt"
	"github.com/shitamachi/promtail-client/promtail"
	"log"
	"os"
	"time"
)

func displayUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s proto|json source-name job-name\n", os.Args[0])
	os.Exit(1)
}

func displayInvalidName(arg string) {
	_, _ = fmt.Fprintf(os.Stderr, "Invalid %s: allowed characters are a-zA-Z0-9_-\n", arg)
	os.Exit(1)
}

func nameIsValid(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			(c == '-') || (c == '_')) {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 4 {
		displayUsage()
	}

	format := os.Args[1]
	sourceName := os.Args[2]
	jobName := os.Args[3]
	if format != "proto" && format != "json" {
		displayUsage()
	}

	if !nameIsValid(sourceName) {
		displayInvalidName("source-name")
	}

	if !nameIsValid(jobName) {
		displayInvalidName("job-name")
	}

	labels := "{source=\"" + sourceName + "\",job=\"" + jobName + "\"}"
	conf := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	var (
		loki promtail.Client
		err  error
	)

	if format == "proto" {
		loki, err = promtail.NewClientProto(conf)
	} else {
		loki, err = promtail.NewClientJson(conf)
	}

	if err != nil {
		log.Printf("promtail.NewClient: %s\n", err)
		os.Exit(1)
	}

	for i := 1; i < 5; i++ {
		tstamp := time.Now().String()
		loki.Debugf("source = %s time = %s, i = %d\n", sourceName, tstamp, i)
		loki.Infof("source = %s, time = %s, i = %d\n", sourceName, tstamp, i)
		loki.Warnf("source = %s, time = %s, i = %d\n", sourceName, tstamp, i)
		loki.Errorf("source = %s, time = %s, i = %d\n", sourceName, tstamp, i)
		time.Sleep(1 * time.Second)
	}

	loki.Shutdown()
}
