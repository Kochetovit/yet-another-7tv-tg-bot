package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	tgBotTokenKey      = "TG_BOT_TOKEN"
	tgUserIDKey        = "TG_USER_ID"
	tgBotNameKey       = "TG_BOT_NAME"
	stickerPackNameKey = "STICKER_PACK_NAME"
	chatIDKey          = "CHAT_ID"

	url7tvKey = "URL_7TV"

	portKey = "PORT"
)

type Config struct {
	// env
	TelegramBotToken string

	UserID  int
	ChatID  int64
	BotName string

	StickerPackName string

	URL7tv string

	Port uint64
}

func Init() (*Config, error) {
	c := &Config{}

	var err error

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

	chatID := os.Getenv(chatIDKey)

	c.ChatID, err = strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		err = errors.Join(
			errors.New(chatIDKey+" is not set"),
			err,
		)

		return nil, err
	}

	c.URL7tv = os.Getenv(url7tvKey)
	if c.URL7tv == "" {
		return nil, errors.New(url7tvKey + " is not set")
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
