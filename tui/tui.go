package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"n1h41/bybit-wallet-tui/repository"
	"n1h41/bybit-wallet-tui/tui/constants"
)

type viewState int

const (
	walletView viewState = iota
	entryView
)

type keyMap struct {
	Quit    key.Binding
	Help    key.Binding
	Balance key.Binding
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
		{k.Balance},
	}
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Show help"),
	),
	Balance: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Get wallet balance"),
	),
}

type mainModel struct {
	state viewState
	keys  keyMap
	help  help.Model
	repo  repository.BybitRepository
}

func NewEntryModel() (tea.Model, tea.Cmd) {
	return mainModel{
		keys: keys,
		help: help.New(),
	}, nil
}

// Init implements tea.Model.
func (m mainModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case "w":
			walletModel := NewWalletModel(constants.Repo)
			return walletModel, walletModel.Init()
		case "d":
			depositModel := NewDepositModel()
			return depositModel, depositModel.Init()
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

// View implements tea.Model.
func (m mainModel) View() string {
	windowSize := constants.WindowSize
	contentView := fmt.Sprintf("%s\n%s", "n1h41", "Bybit Wallet TUI")
	helpView := m.help.View(m.keys)
	helpView = lipgloss.NewStyle().MarginTop(windowSize.Height / 2).Render(helpView)
	return lipgloss.Place(windowSize.Width, windowSize.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, contentView, helpView))
}
