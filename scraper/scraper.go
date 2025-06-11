package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func cleanText(input string) string {
	input = strings.ReplaceAll(input, "\u200b", "")
	input = regexp.MustCompile(`\n{3,}`).ReplaceAllString(input, "\n")
	input = strings.ReplaceAll(input, "\n\n", " ")

	return input
}

func scrapeFirstFoodArticleURL() (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1200, 900),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	var articleURL string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.naver.com"),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`//a[text()="푸드"]`, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.AttributeValue(
			`//a[contains(@class, "TwoColumnView-module__link_item")]`,
			"href",
			&articleURL,
			nil,
		),
	)
	if err != nil {
		log.Println("error finding article URL:", err)
		return "", err
	}

	return articleURL, nil
}

func scrapeArticleContent(articleURL string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1200, 900),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
	defer timeoutCancel()

	var iframeSrc string
	var articleBody string

	err := chromedp.Run(ctx,
		// Navigate to main article page
		chromedp.Navigate(articleURL),
		chromedp.Sleep(1*time.Second),

		// Extract iframe's src attribute
		chromedp.AttributeValue(`iframe#mainFrame`, "src", &iframeSrc, nil),
	)
	if err != nil {
		log.Println("Failed to get iframe src:", err)
		return "", err
	}

	fullIframeURL := articleURL // default
	if iframeSrc != "" {
		fullIframeURL = "https://blog.naver.com" + iframeSrc
	}

	err = chromedp.Run(ctx,
		// Navigate to iframe source directly
		chromedp.Navigate(fullIframeURL),
		chromedp.Sleep(1*time.Second),

		// Get visible text inside the se-viewer content div
		chromedp.Text(`div.se-viewer`, &articleBody, chromedp.NodeVisible),
	)
	if err != nil {
		log.Println("Failed to scrape article content:", err)
		return "", err
	}

	cleaned := cleanText(articleBody)
	return cleaned, nil
}

func ScrapeArticle() (string, string) {
	url, err := scrapeFirstFoodArticleURL()
	if err != nil {
		log.Fatal("Failed to get article URL:", err)
	}
	fmt.Println("Found article URL")

	article, err := scrapeArticleContent(url)
	if err != nil {
		log.Fatal("Failed to get article content:", err)
	}
	fmt.Println("Scraped article content")

	return url, article
}
