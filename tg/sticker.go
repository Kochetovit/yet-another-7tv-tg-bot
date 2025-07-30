package tg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"7tv/tg/dto"
)

type StickerType int

const (
	stickerTypePNG StickerType = iota
	stickerTypeGIF
)

type size int

const (
	SizeS size = 128
	SizeM size = 256
	SizeL size = 512
)

type addStickerOptions struct {
	size   size
	emotes []string
}

func newAddStickerOptions() *addStickerOptions {
	return &addStickerOptions{
		size:   SizeL,
		emotes: []string{"💀"},
	}
}

type addStickerOptionsFn func(opts *addStickerOptions)

func WithSize(size size) addStickerOptionsFn {
	return func(opts *addStickerOptions) {
		opts.size = size
	}
}

func WithEmotes(emotes []string) addStickerOptionsFn {
	return func(opts *addStickerOptions) {
		opts.emotes = emotes
	}
}

func (t *Bot) AddStickerToSet(
	url string,
	stickerType StickerType,
	opts ...addStickerOptionsFn,
) error {
	slog.Info("start adding sticker", "url", url, "stickerType", stickerType)

	o := newAddStickerOptions()
	for _, fn := range opts {
		fn(o)
	}

	emoteID := path.Base(url)

	var (
		suffix    string
		format    dto.StickerFormat
		convertFn func(string, size) (string, error)
	)

	switch stickerType {
	case stickerTypePNG:
		suffix = pngSuffix
		format = dto.StickerFormatStatic
		convertFn = convertPNG
	case stickerTypeGIF:
		suffix = gifSuffix
		format = dto.StickerFormatVideo
		convertFn = convertGIF
	}

	tmpInput, err := downloadImage(t.cfg.URL7tv + "/" + emoteID + "/" + suffix)
	if err != nil {
		slog.Warn("failed to download file")

		return err
	}
	defer os.Remove(tmpInput)

	slog.Info("downloaded sticker", "url", url, "stickerType", stickerType)

	tmpOutput, err := convertFn(tmpInput, o.size)
	if err != nil {
		slog.Warn("failed to convert file")

		return err
	}
	defer os.Remove(tmpOutput)

	slog.Info("converted sticker", "url", url, "stickerType", stickerType)

	u := t.constructURL("/addStickerToSet")

	fields, err := dto.MarshalMap(
		dto.AddStickerToSetRequest{
			UserID: t.cfg.UserID,
			Name:   t.cfg.StickerPackName + "_by_" + t.cfg.BotName,
			Sticker: dto.InputSticker{
				Sticker:   attachKey + fileKey,
				Format:    format,
				EmojiList: o.emotes,
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

	if res.StatusCode != http.StatusOK {
		err = errors.New("failed to add sticker")

		slog.Error("status not ok", "err", err.Error())

		return err
	}

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

	if res.StatusCode != http.StatusOK {
		err = errors.New("failed to create sticker set")
		slog.Error("status not ok", "err", err.Error())

		return err
	}

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
