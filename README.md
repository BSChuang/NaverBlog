# NaverBlog Project

A Go-based project that integrates with Naver Blog to scrape food articles and create interactive quiz content, with Discord integration capabilities.

## Project Structure

- `main.go` - Main entry point of the application
- `scraper/` - Web scraping functionality for Naver Blog articles
- `quiz/` - Quiz generation and management
- `discord/` - Discord bot integration

## Features

- Web scraping of Naver Blog food articles using ChromeDP
- Quiz generation from scraped content
- Discord bot integration for interactive quiz delivery

## Prerequisites

- Go 1.16 or higher
- Chrome/Chromium browser installed
- Discord bot token (for Discord integration)

## Installation

1. Clone the repository:
```bash
git clone [repository-url]
cd NaverBlog
```

2. Install dependencies:
```bash
go get github.com/chromedp/chromedp

```

## Usage

Run the main application:
```bash
go run main.go
```

## Components

### Scraper
The scraper component uses ChromeDP to navigate Naver Blog and extract food article URLs. It's configured to run in non-headless mode for debugging purposes.

### Quiz
The quiz component generates interactive quizzes based on the scraped content.

### Discord
The Discord component provides bot integration for delivering quizzes and managing user interactions.

## Development

To contribute to this project:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 