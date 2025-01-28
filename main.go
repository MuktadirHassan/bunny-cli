package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/MuktadirHassan/bunny-cli/bunny"
	"github.com/lmittmann/tint"
)

func main() {
	// create a new logger
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    false,
		}),
	))
	err := bunny.UploadFolder("./out", 5)
	if err != nil {
		slog.Debug(err.Error())
	}
}
