package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/MuktadirHassan/bunny-cli/cmd"
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

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
