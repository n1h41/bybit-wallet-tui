package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
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
		key.WithKeys("w"),
		key.WithHelp("w", "Open wallet"),
	),
}

type mainModel struct {
	size    tea.WindowSizeMsg
	spinner spinner.Model
	state   viewState
	keys    keyMap
	help    help.Model
	repo    repository.BybitRepository
}

func NewEntryModel() (tea.Model, tea.Cmd) {
	s := spinner.New()
	s.Spinner = spinner.Globe
	return mainModel{
		spinner: s,
		keys:    keys,
		help:    help.New()
	}, nil
}

// Init implements tea.Model.
func (m mainModel) Init() tea.Cmd {
	log.Println("Init")
	return m.spinner.Tick
}

// Update implements tea.Model.
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var sCmd tea.Cmd
	m.spinner, sCmd = m.spinner.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Println(msg)
		m.size = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case "w":
			walletModel := NewWalletModel(constants.Repo, m.size)
			return walletModel, tea.Batch(walletModel.Init())
		case "d":
			depositModel := NewDepositModel(m.size)
			return depositModel, depositModel.Init()
		default:
			return m, nil
		}
	default:
		return m, tea.Batch(sCmd)
	}
}

// View implements tea.Model.
func (m mainModel) View() string {
	windowSize := constants.WindowSize
	contentView := fmt.Sprintf("%s\n%s", "n1h41", "Bybit Wallet")
	helpView := m.help.View(m.keys)
	helpView = lipgloss.NewStyle().MarginTop(windowSize.Height / 2).Render(helpView)
	return lipgloss.Place(windowSize.Width, windowSize.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, contentView, helpView, m.spinner.View()))
}
