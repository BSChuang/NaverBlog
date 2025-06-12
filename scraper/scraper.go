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

func chromeContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-zygote", true),
		chromedp.WindowSize(1200, 900),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, timeout)
	cancel := func() {
		timeoutCancel()
		cancelCtx()
		cancelAlloc()
	}
	return timeoutCtx, cancel
}

func scrapeFirstFoodArticleURL() (string, error) {
	ctx, cancel := chromeContext(15 * time.Second)
	defer cancel()

	var articleURL string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.naver.com"),
		chromedp.Sleep(1*time.Second),
		chromedp.WaitReady("body", chromedp.ByQuery),
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
	ctx, cancel := chromeContext(30 * time.Second)
	defer cancel()

	var iframeSrc string
	var articleBody string

	err := chromedp.Run(ctx,
		chromedp.Navigate(articleURL),
		chromedp.Sleep(1*time.Second),
		chromedp.AttributeValue(`iframe#mainFrame`, "src", &iframeSrc, nil),
	)
	if err != nil {
		log.Println("Failed to get iframe src:", err)
		return "", err
	}

	fullIframeURL := articleURL
	if iframeSrc != "" {
		fullIframeURL = "https://blog.naver.com" + iframeSrc
	}

	err = chromedp.Run(ctx,
		chromedp.Navigate(fullIframeURL),
		chromedp.Sleep(1*time.Second),
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
