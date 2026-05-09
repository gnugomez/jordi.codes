package components

import (
	"strings"

	"jordi.codes/tui/layout"
)

var asciiCat = []string{
	`  ᶻ 𝗓 𐰁 .ᐟ`,
	`/ᐠ > ˕ <マ `,
}

type menuLeftPanel struct {
	cursor int
	width  int
}

func (p menuLeftPanel) render(m AppContext, height int) string {
	var rows []string
	for i, entry := range m.Menu() {
		if i == p.cursor {
			rows = append(rows, selectedItemStyle.Render(entry.Label))
		} else {
			rows = append(rows, normalItemStyle.Render(entry.Label))
		}
	}
	subtitle := mutedStyle.Render("◇") + subtitleStyle.Render("  "+m.Subtitle()+"  ") + mutedStyle.Render("◇")
	block := renderHeader() + subtitle + "\n\n" + strings.Join(rows, "\n")
	return layout.PlaceCentered(p.width, height, block)
}

type menuWidgetPanel struct{ width int }

func (p menuWidgetPanel) render(m AppContext, height int) string {
	return RenderMenuWidgetPanel(MenuWidgetPanelParams{
		Width:      p.width,
		Height:     height,
		Now:        m.Now(),
		Contribs:   m.Contribs(),
		RemoteAddr: m.RemoteAddr(),
	})
}

func renderHeader() string {
	var sb strings.Builder
	for _, line := range asciiCat {
		sb.WriteString(asciiGlowStyle.Render(line))
		sb.WriteString("\n")
	}
	sb.WriteString(titleStyle.Render("Jordi Gómez Hidalgo"))
	sb.WriteString("\n")
	return sb.String()
}
