package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/current/battery-mon/battery"
)

// Styles
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	danger    = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F55385"}

	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2).
			Height(8).
			Width(30).
			Padding(1)

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2).
			Width(60)

	titleStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			MarginBottom(1)

	// Sparkline-ish implementation (simple bar for now)
	barStyle = lipgloss.NewStyle().
			Foreground(special)
)

type tickMsg time.Time
type advTickMsg time.Time

type model struct {
	basic    battery.BatteryInfo
	advanced battery.AdvancedInfo
	err      error
	width    int
	height   int
}

func initialModel() model {
	b, _ := battery.GetBasicInfo()
	// Advanced info takes longer, load it async or just init empty
	return model{
		basic: b,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		advTickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "r" {
			// Force refresh
			b, err := battery.GetBasicInfo()
			m.basic = b
			m.err = err
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		b, err := battery.GetBasicInfo()
		m.basic = b
		m.err = err
		return m, tickCmd()

	case advTickMsg:
		a, err := battery.GetAdvancedInfo()
		if err == nil {
			m.advanced = a
		}
		return m, advTickCmd()
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}

	// Status Color
	statusColor := special
	if m.basic.Percent < 20 {
		statusColor = danger
	}

	// Big Percentage
	percentStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(26).
		Height(3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle)

	pctView := percentStyle.Render(fmt.Sprintf("\n%d%%", m.basic.Percent))

	// Status text
	statusText := fmt.Sprintf("Status: %s\nTime:   %s", m.basic.Status, m.basic.Remaining)

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Battery Monitor"),
		pctView,
		lipgloss.NewStyle().MarginTop(1).Render(statusText),
	)

	// Advanced Stats
	advText := fmt.Sprintf(
		"Condition: %s\nCycles:    %d\nMax Cap:   %s\n\nCharger:   %s\nWattage:   %s\nSerial:    %s",
		m.advanced.Condition,
		m.advanced.CycleCount,
		m.advanced.MaxCapacity,
		m.advanced.ChargerName,
		m.advanced.Wattage,
		m.advanced.Serial,
	)

	rightCol := detailStyle.Render(
		titleStyle.Render("Detailed Stats") + "\n" + advText,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(leftCol), rightCol),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func advTickCmd() tea.Cmd {
	// Poll advanced stats less frequently (every 60s)
	return tea.Tick(time.Second*60, func(t time.Time) tea.Msg {
		return advTickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
