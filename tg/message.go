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

	route := parts[0]

	var cmd command
	switch route {
	case "/start":
		cmd = &startCommand{}
	case "/add":
		cmd = &addCommand{}
	default:
		cmd = &helpCommand{}
	}

	if err := cmd.parseFlags(parts); err != nil {
		return err
	}

	// move to exec
	responseText, err := cmd.exec(t)
	if err != nil {
		return err
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
