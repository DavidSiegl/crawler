package main

import (
	"encoding/json"
	"os"
	"sort"
)

// writeJSONReport marshals the crawled pages to JSON (sorted by normalized URL
// for deterministic output) and writes the result to filename.
func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	report := make([]PageData, 0, len(pages))
	for _, key := range keys {
		report = append(report, pages[key])
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
