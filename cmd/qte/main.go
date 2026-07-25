package main

import (
	"log"
	"log/slog"
	"os"
	"query-terminal/internal/root"

	tea "charm.land/bubbletea/v2"
)

func main() {
	setupLogger()

	p := tea.NewProgram(root.NewRootModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func setupLogger() {
	file, err := os.OpenFile("qte.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error encountered in the server", slog.Any("error message", err))
		os.Exit(1)
	}

	jsonHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(jsonHandler))
}
