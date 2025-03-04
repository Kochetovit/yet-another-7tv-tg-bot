package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"7tv/tg/dto"
)

func (t *Bot) Run() error {
	if err := t.setWebhook(); err != nil {
		slog.Error("failed to set webhook", "err", err.Error())

		return err
	}

	http.Handle("/", http.HandlerFunc(t.UpdateHandler))

	return http.ListenAndServe(":"+strconv.Itoa(int(t.cfg.Port)), nil)
}

func (t *Bot) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
	}()

	var update dto.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		slog.Error("Error decoding update", "err", err.Error())

		return
	}

	go func() {
		if err := t.processUpdate(update); err != nil {
			slog.Error("Error processing update", "err", err.Error())
		}
	}()
}

func (t *Bot) setWebhook() error {
	u := t.constructURL("/setWebhook")

	reqBody, err := json.Marshal(
		dto.SetWebhook{
			URL: t.cfg.WebhookURL,
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
