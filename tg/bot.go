package tg

import (
	"net/http"
	"time"

	"7tv/tg/config"
)

const (
	telegramAPIBaseURL = "https://api.telegram.org/bot"

	fileKey   = "file"
	attachKey = "attach://"
)

type Bot struct {
	cfg    *config.Config
	client *http.Client
}

func NewBot(cfg *config.Config) *Bot {
	return &Bot{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *Bot) constructURL(path string) string {
	return telegramAPIBaseURL + t.cfg.TelegramBotToken + path
}
