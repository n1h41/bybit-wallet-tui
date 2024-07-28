package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"n1h41/bybit-wallet-tui/tui/constants"
)

type depositModel struct {
	size tea.WindowSizeMsg
}

func NewDepositModel(size tea.WindowSizeMsg) tea.Model {
	return depositModel{
		size: size,
	}
}

// Init implements tea.Model.
func (d depositModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (d depositModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		return d, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			mainModel, _ := NewEntryModel(d.size)
			return mainModel, mainModel.Init()
		case "q":
			return d, tea.Quit
		default:
			return d, nil
		}
	default:
		return d, nil
	}
}

// View implements tea.Model.
func (d depositModel) View() string {
	return "Deposit View"
}
