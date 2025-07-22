package tg

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"7tv/tg/dto"
)

func (t *Bot) Run() error {
	http.Handle("/update", http.HandlerFunc(t.UpdateHandler))

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

	if err := t.processUpdate(update); err != nil {
		slog.Error("Error processing update", "err", err.Error())
	}
}
