package main

import (
	"context"
	"fmt"
	"os"

	"keeper/internal/infra/api/grpc"
	"keeper/internal/infra/tui"
	"keeper/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	buildDate    = "N/A"
	buildVersion = "N/A"
)

func main() {
	ctx := context.Background()
	logger.Initialize("debug")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	client := grpc.NewClient(ctx, "127.0.0.1:8080")
	defer client.Close()
	if _, err := tea.NewProgram(tui.InitModel(client, buildDate, buildVersion), tea.WithAltScreen(), tea.WithContext(ctx)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
