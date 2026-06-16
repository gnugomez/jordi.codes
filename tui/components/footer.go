package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Footer struct{}

func (Footer) Render(m AppContext) string {
	t := m.Theme()
	separator := t.MutedStyle.Render("◇ " + strings.Repeat("─", max(m.Width()-2, 0)))
	clockStr := m.Now().Format("Mon 02 Jan 2006  15:04:05")
	clockView := t.ClockStyle.Render(clockStr)
	helpView := t.HelpStyle.Render(m.HelpText())

	lineWidth := m.Width() - 2
	if lineWidth < 1 {
		lineWidth = 1
	}

	clockWidth := lipgloss.Width(clockView)
	helpWidth := lineWidth - clockWidth
	if helpWidth < 1 {
		helpWidth = 1
	}
	helpSlot := t.NewStyle().Width(helpWidth).Render(helpView)
	line := " " + lipgloss.JoinHorizontal(lipgloss.Top, helpSlot, clockView) + " "
	return separator + "\n" + line
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
