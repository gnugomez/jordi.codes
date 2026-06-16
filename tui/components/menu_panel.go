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
	t := m.Theme()
	var rows []string
	for i, entry := range m.Menu() {
		if i == p.cursor {
			rows = append(rows, t.SelectedItemStyle.Render(entry.Label))
		} else {
			rows = append(rows, t.NormalItemStyle.Render(entry.Label))
		}
	}
	subtitle := t.MutedStyle.Render("◇") + t.SubtitleStyle.Render("  "+m.Subtitle()+"  ") + t.MutedStyle.Render("◇")
	block := renderHeader(t) + subtitle + "\n\n" + strings.Join(rows, "\n")
	return layout.PlaceCentered(p.width, height, block)
}

type menuWidgetPanel struct{ width int }

func (p menuWidgetPanel) render(m AppContext, height int) string {
	return MenuWidgetPanel{Width: p.width, Height: height}.Render(m)
}

func renderHeader(t *Theme) string {
	var sb strings.Builder
	for _, line := range asciiCat {
		sb.WriteString(t.AsciiGlowStyle.Render(line))
		sb.WriteString("\n")
	}
	sb.WriteString(t.TitleStyle.Render("Jordi Gómez Hidalgo"))
	sb.WriteString("\n")
	return sb.String()
}
