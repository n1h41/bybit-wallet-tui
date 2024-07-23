package constants

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"n1h41/bybit-wallet-tui/repository"
)

var WindowSize tea.WindowSizeMsg

var Repo repository.BybitRepository

type WalletBalanceMsg struct {
	Rows []table.Row
}
