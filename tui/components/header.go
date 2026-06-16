package components

import (
	"strings"
)

type SectionHeader struct{ Title string }

func (h SectionHeader) Render(m AppContext) string {
	t := m.Theme()
	titleLine := t.NewStyle().
		PaddingLeft(2).
		Render(t.MutedStyle.Render("◈  ") + t.HeaderStyle.Render(h.Title))
	separator := t.MutedStyle.Render(strings.Repeat("─", m.Width()))
	return titleLine + "\n" + separator
}
