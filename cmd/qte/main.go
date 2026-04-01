package main

import (
	"fmt"
	"log/slog"
	"os"
	"qte/internal/db"
	"qte/internal/root"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("this is qte")
	setupLogger()
	db.InitializeDBRegistry()

	slog.Debug("starting application...")
	p := tea.NewProgram(root.NewRootModel())
	if _, err := p.Run(); err != nil {
		slog.Error("Error:", err)
		os.Exit(1)
	}
}

func setupLogger() {
	file, err := os.OpenFile("qte.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error:", err)
		os.Exit(1)
	}

	jsonHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(jsonHandler))
}
