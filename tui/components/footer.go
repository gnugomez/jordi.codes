package components

import "strings"

type Footer struct{}

func (Footer) Render(m AppContext) string {
	separator := mutedStyle.Render("◇ " + strings.Repeat("─", m.Width()-2))
	clockStr := m.Now().Format("Mon 02 Jan 2006  15:04:05")
	clockView := clockStyle.Render(clockStr)
	helpView := helpStyle.Render(m.HelpText())
	gap := m.Width() - width(helpView) - width(clockView) - 2
	if gap < 1 {
		gap = 1
	}
	line := " " + helpView + strings.Repeat(" ", gap) + clockView + " "
	return separator + "\n" + line
}
