package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(htmlBody string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return "", err
	}

	if h1 := doc.Find("h1").First(); h1.Length() > 0 {
		return strings.TrimSpace(h1.Text()), nil
	}
	return strings.TrimSpace(doc.Find("h2").First().Text()), nil
}

func getFirstParagraphFromHTML(htmlBody string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return "", err
	}

	if main := doc.Find("main").First(); main.Length() > 0 {
		if p := main.Find("p").First(); p.Length() > 0 {
			return strings.TrimSpace(p.Text()), nil
		}
	}

	return strings.TrimSpace(doc.Find("p").First().Text()), nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		ref, err := url.Parse(href)
		if err != nil {
			return
		}
		urls = append(urls, baseURL.ResolveReference(ref).String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		ref, err := url.Parse(src)
		if err != nil {
			return
		}
		images = append(images, baseURL.ResolveReference(ref).String())
	})

	return images, nil
}
