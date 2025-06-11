package main

import (
	sender "NaverBlog/discord"
	generator "NaverBlog/openai"
	"NaverBlog/scraper"
)

func main() {
	url, article := scraper.ScrapeArticle()
	cleanedArticle, quiz := generator.CreateQuiz(article)
	sender.SendArticleAndQuiz(url, cleanedArticle, quiz)
}
