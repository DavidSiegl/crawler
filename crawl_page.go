package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// addPageVisit records a visit to normalizedURL. It returns isFirst=true only
// for the first goroutine to reach a given URL, so exactly one of them goes on
// to fetch and crawl the page. It is safe for concurrent use.
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}
	// Reserve the slot so concurrent crawlers see it as already visited.
	cfg.pages[normalizedURL] = PageData{}
	return true
}

// reachedMaxPages reports whether we've already crawled at least maxPages
// pages. It is safe for concurrent use.
func (cfg *config) reachedMaxPages() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) >= cfg.maxPages
}

// crawlPage crawls rawCurrentURL and recursively crawls every internal link it
// finds, bounded by cfg.concurrencyControl. Each call runs in its own goroutine
// and signals completion via cfg.wg.
func (cfg *config) crawlPage(rawCurrentURL string) {
	// Acquire a concurrency slot for the duration of this call.
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// Stop once we've hit the page limit. Checked under the mutex.
	if cfg.reachedMaxPages() {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing current URL %q: %v\n", rawCurrentURL, err)
		return
	}

	// Only crawl pages on the same domain as the base URL.
	if currentURL.Host != cfg.baseURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL %q: %v\n", rawCurrentURL, err)
		return
	}

	// Only the first goroutine to reach this page proceeds to fetch it.
	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML from %q: %v\n", rawCurrentURL, err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
}
