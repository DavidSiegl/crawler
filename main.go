package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("usage: crawler BASE_URL MAX_CONCURRENCY MAX_PAGES")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[0]

	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil || maxConcurrency < 1 {
		fmt.Printf("invalid maxConcurrency %q: must be a positive integer\n", args[1])
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(args[2])
	if err != nil || maxPages < 1 {
		fmt.Printf("invalid maxPages %q: must be a positive integer\n", args[2])
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL %q: %v\n", rawBaseURL, err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              map[string]PageData{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	const reportFile = "report.json"
	if err := writeJSONReport(cfg.pages, reportFile); err != nil {
		fmt.Printf("error writing report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("wrote report for %d pages to %s\n", len(cfg.pages), reportFile)
}
