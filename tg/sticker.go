package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"7tv/tg/dto"
)

func (t *Bot) AddStickerToSet(url string) error {
	emoteID := path.Base(url)

	tmpInput, err := downloadImage(url7tv + emoteID + pngSuffix)
	if err != nil {
		slog.Error("failed to download image", "err", err.Error())

		return err
	}
	defer os.Remove(tmpInput)

	tmpOutput, err := convertPNG(tmpInput)
	if err != nil {
		slog.Error("failed to convert image", "err", err.Error())

		return err
	}
	defer os.Remove(tmpOutput)

	u := t.constructURL("/addStickerToSet")

	fields, err := dto.MarshalMap(
		dto.AddStickerToSetRequest{
			UserID: t.cfg.UserID,
			Name:   t.cfg.StickerPackName + "_by_" + t.cfg.BotName,
			Sticker: dto.InputSticker{
				Sticker:   attachKey + fileKey,
				Format:    dto.StickerFormatStatic,
				EmojiList: []string{"💀"},
			},
		},
	)
	if err != nil {
		slog.Error("failed to marshal fields", "err", err.Error())

		return err
	}

	reqBody := new(bytes.Buffer)
	mw := multipart.NewWriter(reqBody)

	for key, value := range fields {
		if err := mw.WriteField(key, value); err != nil {
			slog.Error("failed to write field", "err", err.Error())

			return err
		}
	}

	pw, err := mw.CreateFormFile(fileKey, tmpOutput)
	if err != nil {
		slog.Error("failed to create form file", "err", err.Error())

		return err
	}

	file, err := os.Open(tmpOutput)
	if err != nil {
		slog.Error("failed to open file", "err", err.Error())

		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		slog.Error("failed to read file", "err", err.Error())

		return err
	}

	if _, err := pw.Write(content); err != nil {
		slog.Error("failed to write file", "err", err.Error())

		return err
	}

	if err := mw.Close(); err != nil {
		slog.Error("failed to close multipart writer", "err", err.Error())

		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, reqBody)
	if err != nil {
		slog.Error("failed to create request", "err", err.Error())

		return err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

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
					Sticker:   "splort.webp",
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
