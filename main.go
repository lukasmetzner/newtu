package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	db, err := OpenDB(cfg.DBFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close() //nolint:errcheck

	m := newModel(cfg, db)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
