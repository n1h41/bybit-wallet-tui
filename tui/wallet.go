package tui

import (
	"log"
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
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

type walletModel struct {
	size                tea.WindowSizeMsg
	spinner             spinner.Model
	loading             bool
	table               table.Model
	repo                repository.BybitRepository
	walletTotalUSDValue float64
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
		log.Println(msg)
		w.size = msg
		return w, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			mainModel, _ := NewEntryModel(w.size)
			return mainModel, mainModel.Init()
		case "r":
			w.loading = true
			return w, tea.Batch(w.spinner.Tick, w.getWalletBalance())
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
	windowSize := w.size
	rows := msg.Rows
	columns := []table.Column{
		{Title: "Coin", Width: 6},
		{Title: "USD Value", Width: 10},
		{Title: "Amount", Width: 10},
	}
	total := calculateTotalOfRows(rows)
	w.walletTotalUSDValue = total
	ts := table.DefaultStyles()
	ts.Header = table.DefaultStyles().Header.Foreground(lipgloss.Color("99"))
	w.table = table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(windowSize.Height-5), table.WithFocused(true), table.WithStyles(ts))
}

func (w walletModel) getWalletBalance() tea.Cmd {
	return func() tea.Msg {
		var rows []table.Row
		result := w.repo.GetWalletBalance()
		for _, item := range result.Result.List[0].Coin {
			rows = append(rows, table.Row{
				item.Coin,
				item.UsdValue,
				item.WalletBalance,
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
	renderer := constants.Renderer
	windowSize := w.size
	totalStr := "Total USD Value: $" + strconv.FormatFloat(w.walletTotalUSDValue, 'f', 2, 64)
	totalStr = renderer.NewStyle().Padding(0, 0, 1, 1).Bold(true).Foreground(lipgloss.Color("200")).Render(totalStr)
	if w.loading {
		return renderer.Place(windowSize.Width, windowSize.Height, lipgloss.Center, lipgloss.Center, w.spinner.View())
	}
	return renderer.Place(windowSize.Width, windowSize.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, baseStyle.Render(w.table.View()), totalStr))
}

func NewWalletModel(repo repository.BybitRepository, size tea.WindowSizeMsg) tea.Model {
	s := spinner.New()
	return walletModel{
		size:    size,
		spinner: s,
		loading: true,
		repo:    repo,
	}
}

func calculateTotalOfRows(rows []table.Row) float64 {
	var total float64
	for _, row := range rows {
		value, _ := strconv.ParseFloat(row[1], 64)
		total += value
	}
	return total
}
