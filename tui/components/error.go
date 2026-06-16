package components

import "jordi.codes/tui/layout"

type ErrorView struct{ ErrText string }

func (e ErrorView) Render(m AppContext) string {
	t := m.Theme()
	bodyHeight := m.Height() - layout.FooterHeight
	face := t.MutedStyle.Render("( •_• )")
	msg := t.ErrorStyle.Render(e.ErrText)
	hint := t.MutedStyle.Render("press esc to go back")
	return layout.PlaceCentered(m.Width(), bodyHeight, face+"\n\n"+msg+"\n"+hint)
}
