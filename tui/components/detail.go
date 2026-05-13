package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// ContentBody is a value-type Renderer that displays rendered markdown content
// in a viewport, with a section header above it.
type ContentBody struct {
	Title    string
	Viewport viewport.Model
}

func (c ContentBody) Render(m AppContext) string {
	head := SectionHeader{Title: c.Title}.Render(m)
	return lipgloss.JoinVertical(lipgloss.Left, head, c.Viewport.View())
}
