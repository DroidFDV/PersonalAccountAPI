package main

import (
	"PersonalAccountAPI/internal/app"
	"log/slog"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1)
	}
}
