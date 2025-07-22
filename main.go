package main

import (
	"log/slog"

	"7tv/tg"
	"7tv/tg/config"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		slog.Error("failed to init config", "err", err.Error())

		return
	}

	tgBot := tg.NewBot(cfg)

	if err := tgBot.Run(); err != nil {
		slog.Error("failed to run bot", "err", err.Error())

		return
	}

	// if err := tgBot.AddStickerToSet("https://7tv.app/emotes/01FFFPWV180007P57XYW0BHF1Z"); err != nil {
	// 	slog.Error("failed to get sticker set", "err", err.Error())
	// }
}
