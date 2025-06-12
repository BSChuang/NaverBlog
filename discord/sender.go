package sender

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func initEnv() (string, string) {
	_ = godotenv.Load()

	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	if discordToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set")
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		log.Fatal("DISCORD_CHANNEL_ID is not set")
	}

	return discordToken, channelID
}

func sendMessage(dg *discordgo.Session, channelID string, message string) error {
	chunks := splitMessage(message, 1900)

	for _, chunk := range chunks {
		_, err := dg.ChannelMessageSend(channelID, chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitMessage(msg string, maxLen int) []string {
	var chunks []string
	runes := []rune(msg) // handle Unicode properly
	for len(runes) > 0 {
		end := maxLen
		if len(runes) < end {
			end = len(runes)
		} else {
			// Try to split on a newline or space for cleaner output
			safe := end
			for i := end - 1; i > end-100 && i > 0; i-- {
				if runes[i] == '\n' || runes[i] == ' ' {
					safe = i
					break
				}
			}
			end = safe
		}
		chunks = append(chunks, string(runes[:end]))
		runes = runes[end:]
	}
	return chunks
}

func SendArticleAndQuiz(url string, article string, quiz string) {
	discordToken, channelID := initEnv()

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("failed to create Discord session: %w", err)
	}
	defer dg.Close()

	if err := dg.Open(); err != nil {
		log.Fatal("failed to open connection to Discord: %w", err)
	}

	err = sendMessage(dg, channelID, fmt.Sprintf("# %s Blog", time.Now().Format("2006-01-02")))
	if err != nil {
		log.Fatal("failed to send time: %w", err)
	}

	err = sendMessage(dg, channelID, url)
	if err != nil {
		log.Fatal("failed to send url: %w", err)
	}

	err = sendMessage(dg, channelID, article)
	if err != nil {
		log.Fatal("failed to send article: %w", err)
	}

	err = sendMessage(dg, channelID, fmt.Sprintf("\n\n--------------------------------\n\n%s", quiz))
	if err != nil {
		log.Fatal("failed to send quiz: %w", err)
	}
	fmt.Println("Sent messages to channel")
}
