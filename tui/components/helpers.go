package components

import "github.com/charmbracelet/lipgloss"

func width(s string) int {
	return lipgloss.Width(s)
}
