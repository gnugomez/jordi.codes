package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type DetailBody struct {
	Title    string
	Viewport viewport.Model
}

func (d DetailBody) Render(m AppContext) string {
	head := SectionHeader{Title: d.Title}.Render(m)
	return lipgloss.JoinVertical(lipgloss.Left, head, d.Viewport.View())
}
