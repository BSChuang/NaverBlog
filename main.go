package main

import (
	generator "NaverBlog/openai"
	"NaverBlog/scraper"
	"fmt"
)

func main() {
	article := scraper.ScrapeArticle()
	cleanedArticle, quiz := generator.CreateQuiz(article)
	fmt.Print(cleanedArticle, quiz)
}
