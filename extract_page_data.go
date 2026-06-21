package main

import "net/url"

// PageData holds the structured data extracted from a single crawled page.
type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

// extractPageData parses an HTML page and returns its structured data. pageURL
// is the absolute URL of the page and is used to resolve any relative links and
// image sources into absolute URLs.
func extractPageData(html, pageURL string) PageData {
	page := PageData{URL: pageURL}

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return page
	}

	if heading, err := getHeadingFromHTML(html); err == nil {
		page.Heading = heading
	}
	if paragraph, err := getFirstParagraphFromHTML(html); err == nil {
		page.FirstParagraph = paragraph
	}
	if links, err := getURLsFromHTML(html, baseURL); err == nil {
		page.OutgoingLinks = links
	}
	if images, err := getImagesFromHTML(html, baseURL); err == nil {
		page.ImageURLs = images
	}

	return page
}
