package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	tgWebhookURLKey    = "TG_WEBHOOK_URL"
	tgBotTokenKey      = "TG_BOT_TOKEN"
	tgUserIDKey        = "TG_USER_ID"
	tgBotNameKey       = "TG_BOT_NAME"
	stickerPackNameKey = "STICKER_PACK_NAME"
	portKey            = "PORT"
)

type Config struct {
	// env
	TelegramBotToken string
	WebhookURL       string

	UserID  int
	BotName string

	StickerPackName string

	Port uint64
}

func Init() (*Config, error) {
	c := &Config{}

	var err error

	c.WebhookURL = os.Getenv(tgWebhookURLKey)
	if c.WebhookURL == "" {
		return nil, errors.New(tgWebhookURLKey + " is not set")
	}

	c.TelegramBotToken = os.Getenv(tgBotTokenKey)
	if c.TelegramBotToken == "" {
		return nil, errors.New(tgBotTokenKey + " is not set")
	}

	userID := os.Getenv(tgUserIDKey)

	c.UserID, err = strconv.Atoi(userID)
	if err != nil {
		err = errors.Join(
			errors.New(tgUserIDKey+" is not set"),
			err,
		)

		return nil, err
	}

	c.BotName = os.Getenv(tgBotNameKey)
	if c.BotName == "" {
		return nil, errors.New(tgBotNameKey + " is not set")
	}

	c.StickerPackName = os.Getenv(stickerPackNameKey)
	if c.StickerPackName == "" {
		return nil, errors.New(stickerPackNameKey + " is not set")
	}

	port := os.Getenv(portKey)

	c.Port, err = strconv.ParseUint(port, 0, 64)
	if err != nil {
		err = errors.Join(
			errors.New(portKey+" is not set"),
			err,
		)

		return nil, err
	}

	return c, nil
}
