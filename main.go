package main

import (
	"log/slog"

	"7tv/tg/config"
	"7tv/tg"
)

func main() {
	cfg, err := config.Init("config.json")
	if err != nil {
		slog.Error("failed to init config", "err", err.Error())

		return
	}

	tgBot := tg.NewBot(cfg)

	if err := tgBot.AddFileStickerToSet("https://7tv.app/emotes/01FFFPWV180007P57XYW0BHF1Z"); err != nil {
		slog.Error("failed to get sticker set", "err", err.Error())
	}
}
