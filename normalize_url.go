package main

import (
	"net/url"
	"strings"
)

// normalizeURL returns a canonical form of raw URL string suitable for
// deduplication: the scheme is dropped and the host is lowercased, leaving
// "host/path" with any trailing slash removed.
func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}

	normalized := strings.ToLower(u.Host) + u.Path
	normalized = strings.TrimSuffix(normalized, "/")

	return normalized, nil
}
