package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

const (
	tgBotTokenKey      = "TG_BOT_TOKEN"
	tgUserIDKey        = "TG_USER_ID"
	tgBotNameKey       = "TG_BOT_NAME"
	stickerPackNameKey = "STICKER_PACK_NAME"
)

type Config struct {
	// env
	TelegramBotToken string
	UserID           int
	BotName          string
	StickerPackName  string

	// config
	TelegramAPIBaseURL string `json:"telegramApiBaseUrl"`
	Splort             string `json:"splort"`
}

func Init(configPath string) (*Config, error) {
	c := &Config{}

	var (
		b   []byte
		err error
	)

	b, err = os.ReadFile(configPath)
	if err != nil {
		err = errors.Join(
			errors.New("failed to read config file"),
			err,
		)

		return nil, err
	}

	if err = json.Unmarshal(b, c); err != nil {
		err = errors.Join(
			errors.New("failed to unmarshal config file"),
			err,
		)

		return nil, err
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

	return c, nil
}
