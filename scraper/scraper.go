package scraper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func ScrapeFirstFoodArticleURL() (string, error) {
	// Use non-headless Chrome for debugging
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // show browser window
		chromedp.Flag("disable-gpu", false),
		chromedp.WindowSize(1200, 900),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Optional: set a timeout to prevent freezing
	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	var articleURL string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.naver.com"),
		chromedp.Sleep(3*time.Second),

		// Click 푸드 category (you may need to update selector)
		chromedp.Click(`//a[text()="푸드"]`, chromedp.BySearch),
		chromedp.Sleep(3*time.Second),

		// Get href of first article link on 푸드 page
		chromedp.AttributeValue(
			`//a[contains(@class, "TwoColumnView-module__link_item")]`,
			"href",
			&articleURL,
			nil,
		),
	)

	if err != nil {
		log.Println("chromedp run error:", err)
		return "", err
	}

	return articleURL, nil
}
