package main

import (
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"n1h41/bybit-wallet-tui/configs"
	"n1h41/bybit-wallet-tui/repository"
	"n1h41/bybit-wallet-tui/tui"
	"n1h41/bybit-wallet-tui/tui/constants"
)

func main() {
	f, _ := tea.LogToFile("debug.log", "Debug: ")
	defer f.Close()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	bybitConfig := configs.BybitConfig{
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		RecvWindow: "5000",
		Url:        "https://api.bybit.com",
	}

	client := &http.Client{Timeout: 10 * time.Second}

	constants.Repo = repository.NewBybitRepo(bybitConfig, client)

	m, _ := tui.NewEntryModel(tea.WindowSizeMsg{})

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
