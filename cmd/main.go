package main

import (
	"PersonalAccountAPI/internal/app"
	"log/slog"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("main app.Run", slog.Any("error", err))
		os.Exit(1)
	}
}
