package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/curtisepp/notion-tui/internal/config"
	"github.com/curtisepp/notion-tui/internal/notion"
	"github.com/curtisepp/notion-tui/internal/tui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("notion-tui", version)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, "Set NOTION_API_TOKEN environment variable or create ~/.notion-tui/config.yaml:")
		fmt.Fprintln(os.Stderr, "  token: ntn_xxxx...")
		os.Exit(1)
	}

	client := notion.NewClient(cfg.Token)
	app := tui.NewApp(client)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
