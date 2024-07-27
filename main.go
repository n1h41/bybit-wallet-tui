package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"

	"n1h41/bybit-wallet-tui/configs"
	"n1h41/bybit-wallet-tui/repository"
	"n1h41/bybit-wallet-tui/tui"
	"n1h41/bybit-wallet-tui/tui/constants"
)

const (
	host = "127.0.0.1"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

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

	constants.Renderer = bubbletea.MakeRenderer(s)
	constants.WindowSize.Width = pty.Window.Width
	constants.WindowSize.Height = pty.Window.Height

	client := &http.Client{Timeout: 10 * time.Second}

	constants.Repo = repository.NewBybitRepo(bybitConfig, client)

	m, _ := tui.NewEntryModel(tea.WindowSizeMsg{})

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
