package main

import (
	sender "github.com/BSChuang/NaverBlog/discord"
	generator "github.com/BSChuang/NaverBlog/openai"
	"github.com/BSChuang/NaverBlog/scraper"
)

func main() {
	url, article := scraper.ScrapeArticle()
	cleanedArticle, quiz := generator.CreateQuiz(article)
	sender.SendArticleAndQuiz(url, cleanedArticle, quiz)
}
