package constants

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"n1h41/bybit-wallet-tui/repository"
)

var Renderer *lipgloss.Renderer

var WindowSize tea.WindowSizeMsg

var Repo repository.BybitRepository

type WalletBalanceMsg struct {
	Rows []table.Row
}
