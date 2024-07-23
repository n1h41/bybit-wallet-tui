package tui

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"n1h41/bybit-wallet-tui/repository"
	"n1h41/bybit-wallet-tui/tui/constants"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type walletModel struct {
	spinner spinner.Model
	loading bool
	table   table.Model
	repo    repository.BybitRepository
}

// Init implements tea.Model.
func (w walletModel) Init() tea.Cmd {
	return tea.Batch(w.spinner.Tick, w.getWalletBalance())
}

// Update implements tea.Model.
func (w walletModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList []tea.Cmd
	var tbCmd, spCmd tea.Cmd
	w.table, tbCmd = w.table.Update(msg)
	w.spinner, spCmd = w.spinner.Update(msg)
	switch msg := msg.(type) {
	case constants.WalletBalanceMsg:
		w.loading = false
		w.generateTable(msg)
		return w, nil
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		return w, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return NewEntryModel()
		case "q":
			return w, tea.Quit
		default:
			return w, nil
		}
	}

	cmdList = append(cmdList, tbCmd, spCmd)
	return w, tea.Batch(cmdList...)
}

func (w *walletModel) generateTable(msg constants.WalletBalanceMsg) {
	rows := msg.Rows
	columns := []table.Column{
		{Title: "Coin", Width: 6},
		{Title: "USD Value", Width: 10},
	}
	total := calculateTotalOfRows(rows)
	rows = append(rows, total)
	rows = append(rows, rows...)
	rows = append(rows, rows...)
	w.table = table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(15), table.WithFocused(true))
}

func (w walletModel) getWalletBalance() tea.Cmd {
	return func() tea.Msg {
		var rows []table.Row
		result := w.repo.GetWalletBalance()
		for _, item := range result.Result.List[0].Coin {
			rows = append(rows, table.Row{
				item.Coin,
				item.UsdValue,
			})
		}
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][1] > rows[j][1]
		})
		return constants.WalletBalanceMsg{Rows: rows}
	}
}

// View implements tea.Model.
func (w walletModel) View() string {
	if w.loading {
		return lipgloss.Place(constants.WindowSize.Width, constants.WindowSize.Height, lipgloss.Center, lipgloss.Center, w.spinner.View())
	}
	return baseStyle.Render(w.table.View())
}

func NewWalletModel(repo repository.BybitRepository) tea.Model {
	s := spinner.New()
	return walletModel{
		spinner: s,
		loading: true,
		repo:    repo,
	}
}

func calculateTotalOfRows(rows []table.Row) table.Row {
	var total float64
	for _, row := range rows {
		value, _ := strconv.ParseFloat(row[1], 64)
		total += value
	}
	return table.Row{"Total", strconv.FormatFloat(total, 'f', 4, 64)}
}
