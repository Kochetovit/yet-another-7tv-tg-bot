package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"7tv/tg/config"
	"7tv/tg/dto"
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
	return t.cfg.TelegramAPIBaseURL + t.cfg.TelegramBotToken + path
}

func (t *Bot) AddStickerToSet(name string) error {
	u := t.constructURL("/addStickerToSet")

	reqBody, err := json.Marshal(
		dto.AddStickerToSetRequest{
			UserID: t.cfg.UserID,
			Name:   name,
			Sticker: dto.InputSticker{
				Sticker:   t.cfg.Splort,
				Format:    dto.StickerFormatStatic,
				EmojiList: []string{"💀"},
			},
		},
	)
	if err != nil {
		slog.Error("failed to marshal request body", "err", err.Error())

		return err
	}

	reader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, u, reader)
	if err != nil {
		slog.Error("failed to create request", "err", err.Error())

		return err
	}

	fmt.Println(string(reqBody))

	// TODO: think about it
	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		slog.Error("failed to send request", "err", err.Error())

		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("failed to read response body", "err", err.Error())

		return err
	}

	fmt.Printf("%+v", string(resBody))

	return nil
}

func (t *Bot) CreateStickerSet(name string) error {
	u := t.constructURL("/createNewStickerSet")

	reqBody, err := json.Marshal(
		dto.CreateNewStickerSet{
			UserID: t.cfg.UserID,
			Name:   name + "_by_" + t.cfg.BotName,
			Title:  name,
			Stickers: []dto.InputSticker{
				{
					Sticker:   t.cfg.Splort,
					Format:    dto.StickerFormatStatic,
					EmojiList: []string{"💀"},
				},
			},
		},
	)
	if err != nil {
		slog.Error("failed to marshal request body", "err", err.Error())

		return err
	}

	reader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, u, reader)
	if err != nil {
		slog.Error("failed to create request", "err", err.Error())

		return err
	}

	fmt.Println(string(reqBody))

	// TODO: think about it
	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		slog.Error("failed to send request", "err", err.Error())

		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("failed to read response body", "err", err.Error())

		return err
	}

	fmt.Printf("%+v", string(resBody))

	return nil
}

func (t *Bot) GetStickerSet(name string) error {
	u := t.constructURL("/getStickerSet")

	reqBody, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: name,
		},
	)
	if err != nil {
		slog.Error("failed to marshal request body", "err", err.Error())

		return err
	}

	reader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, u, reader)
	if err != nil {
		slog.Error("failed to create request", "err", err.Error())

		return err
	}

	fmt.Println(string(reqBody))

	// TODO: think about it
	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		slog.Error("failed to send request", "err", err.Error())

		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("failed to read response body", "err", err.Error())

		return err
	}

	fmt.Printf("%+v", string(resBody))

	return nil
}
