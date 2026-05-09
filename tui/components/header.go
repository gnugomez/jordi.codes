package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SectionHeader struct{ Title string }

func (h SectionHeader) Render(m AppContext) string {
	titleLine := lipgloss.NewStyle().
		PaddingLeft(2).
		Render(mutedStyle.Render("◈  ") + headerStyle.Render(h.Title))
	separator := mutedStyle.Render(strings.Repeat("─", m.Width()))
	return titleLine + "\n" + separator
}
