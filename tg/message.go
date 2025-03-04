package tg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"7tv/tg/dto"
)

func (t *Bot) processUpdate(update dto.Update) error {
	if update.Message == nil || update.Message.Chat == nil || update.Message.Text == "" {
		return errors.New("invalid update")
	}

	text := update.Message.Text
	chatID := update.Message.Chat.ID

	parts := strings.Fields(text)
	if len(parts) == 0 {
		return errors.New("invalid command")
	}

	command := parts[0]
	var responseText string

	switch command {
	case "/start":
		if err := t.CreateStickerSet(t.cfg.StickerPackName + "_by_" + t.cfg.BotName); err != nil {
			responseText = err.Error()
		}
	case "/png":
		if len(parts) < 2 {
			responseText = "Please provide 7tv emote url"
		} else {
			url := parts[1]

			if err := t.AddStickerToSet(url, stickerTypePNG); err != nil {
				responseText = err.Error()
			} else {
				responseText = "Sticker added to set"
			}
		}
	case "/gif":
		if len(parts) < 2 {
			responseText = "Please provide 7tv emote url"
		} else {
			url := parts[1]

			if err := t.AddStickerToSet(url, stickerTypeGIF); err != nil {
				responseText = err.Error()
			} else {
				responseText = "Sticker added to set"
			}
		}
	default:
		return errors.New("unknown command")
	}

	return t.sendMessage(chatID, responseText)
}

func (t *Bot) sendMessage(chatID int64, text string) error {
	u := t.constructURL("/sendMessage")

	reqBody, err := json.Marshal(
		dto.Message{
			ChatID: int(chatID),
			Text:   text,
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
