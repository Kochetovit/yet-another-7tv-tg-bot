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
}
