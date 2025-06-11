package generator

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type Quiz struct {
	Question string
	Options  []string
	Answer   string
}

func initEnv() string {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	return apiKey
}

func generateQuiz(client *openai.Client, article string) (string, error) {
	prompt := fmt.Sprintf(`Create a 5 question quiz in English on the content of the article. Format it for Discord. Provide the answers in a spoiler tag (i.e. ||Insert spoilers here||). The format should be as follows:
Question: ...
Options:
- ...
- ...
- ...
Answer: ...
%s`, article)
	quiz, err := chat(client, prompt, article)

	return quiz, err
}

func cleanArticle(client *openai.Client, article string) (string, error) {
	prompt := fmt.Sprintf(`Please fix up the punctuation and structure for the following text scraped from a Korean food blog. 
Remove any phrases which are scraped from the website rather than the blog post. Keep the article in Korean. Format it for Discord.

article:
%s`, article)
	cleanArticle, err := chat(client, prompt, article)

	return cleanArticle, err
}

func chat(client *openai.Client, prompt string, content string) (string, error) {
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	})

	if err != nil {
		return "", err
	}

	respString := resp.Choices[0].Message.Content
	return respString, err
}

func CreateQuiz(article string) (string, string) {
	apiKey := initEnv()
	client := openai.NewClient(apiKey)

	cleanedArticle, err := cleanArticle(client, article)
	if err != nil {
		log.Fatal("Error cleaning article:", err)
	}
	fmt.Println("Cleaned article content")

	quiz, err := generateQuiz(client, cleanedArticle)
	if err != nil {
		log.Fatal("Error generating quiz:", err)
	}
	fmt.Println("Generated quiz")

	return cleanedArticle, quiz
}
