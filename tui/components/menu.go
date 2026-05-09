package components

import (
	"github.com/charmbracelet/lipgloss"

	"jordi.codes/tui/layout"
)

type MenuBody struct{ Cursor int }

func (mb MenuBody) Render(m AppContext) string {
	w := MenuWidth(m.Width())
	h := m.Height() - layout.FooterHeight
	if w < layout.MinSplitWidth {
		return menuLeftPanel{cursor: mb.Cursor, width: w}.render(m, h)
	}
	left := menuLeftPanel{cursor: mb.Cursor, width: layout.LeftPanelWidth}.render(m, h)
	right := menuWidgetPanel{width: w - layout.LeftPanelWidth}.render(m, h)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
