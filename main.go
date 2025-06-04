package main

import (
	"fmt"
	"log"

	"NaverBlog/scraper"
)

func main() {
	url, err := scraper.ScrapeFirstFoodArticleURL()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("First food article URL:", url)
}
