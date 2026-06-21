package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "extract h1 tag",
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "extract h2 tag",
			inputBody: "<html><body><h2>Test Title</h2></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "return empty string",
			inputBody: "<html><body><main><p>Test Title</p></main></body></html>",
			expected:  "",
		},
		{
			name:      "prefer h1 tag",
			inputBody: "<html><body><h1>Test Title</h1><h2>Test Title 2</h2></body></html>",
			expected:  "Test Title",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getHeadingFromHTML(tc.inputBody)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected Text: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name: "extract main paragraph",
			inputBody: `<p>Outside paragraph.</p>
					<main>
						<p>Main paragraph.</p>
					</main>`,
			expected: "Main paragraph.",
		},
		{
			name: "",
			inputBody: `<p>Outside paragraph.</p>
					<main>
						<p>Main paragraph.</p>
						<p>Second paragraph.</p>
					</main>`,
			expected: "Main paragraph.",
		},
		{
			name:      "return empty string",
			inputBody: "<html><body><main><h1>Test Title</h1></main></body></html>",
			expected:  "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getFirstParagraphFromHTML(tc.inputBody)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected Text: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "relative URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/path/one"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com/path/one"},
		},
		{
			name:      "multiple URLs",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com/path/one">One</a><a href="/path/two">Two</a></body></html>`,
			expected:  []string{"https://crawler-test.com/path/one", "https://crawler-test.com/path/two"},
		},
		{
			name:      "no links",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No links here.</p></body></html>`,
			expected:  nil,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URLs: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute image URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "relative image URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "multiple images",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/one.png"><img src="/two.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/one.png", "https://crawler-test.com/two.png"},
		},
		{
			name:      "no images",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No images here.</p></body></html>`,
			expected:  nil,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected images: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
