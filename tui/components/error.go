package components

import "jordi.codes/tui/layout"

type ErrorView struct{ ErrText string }

func (e ErrorView) Render(m AppContext) string {
	bodyHeight := m.Height() - layout.FooterHeight
	face := mutedStyle.Render("( •_• )")
	msg := errorStyle.Render(e.ErrText)
	hint := mutedStyle.Render("press esc to go back")
	return layout.PlaceCentered(m.Width(), bodyHeight, face+"\n\n"+msg+"\n"+hint)
}
