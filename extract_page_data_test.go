package main

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  PageData
	}{
		{
			name:     "full page",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html><body>
				<h1>Test Title</h1>
				<main><p>Main paragraph.</p></main>
				<a href="/link/one">One</a>
				<a href="https://other.com/two">Two</a>
				<img src="/logo.png">
				<img src="https://cdn.boot.dev/banner.png">
			</body></html>`,
			expected: PageData{
				URL:            "https://blog.boot.dev",
				Heading:        "Test Title",
				FirstParagraph: "Main paragraph.",
				OutgoingLinks:  []string{"https://blog.boot.dev/link/one", "https://other.com/two"},
				ImageURLs:      []string{"https://blog.boot.dev/logo.png", "https://cdn.boot.dev/banner.png"},
			},
		},
		{
			name:      "empty page",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body></body></html>`,
			expected: PageData{
				URL: "https://blog.boot.dev",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.inputBody, tc.inputURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected: %+v, actual: %+v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
