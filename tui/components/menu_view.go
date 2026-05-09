package components

import (
	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/tui/layout"
)

type MenuView struct {
	cursor int
}

func NewMenuView() *MenuView { return &MenuView{} }

func (mv *MenuView) Render(m AppContext) string {
	m.SetHelpText("↑/↓  navigate   enter  open   q  quit")
	return MenuBody{Cursor: mv.cursor}.Render(m)
}

func (mv *MenuView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch keyMsg.String() {
	case "q":
		return tea.Quit
	case "up", "k":
		if mv.cursor > 0 {
			mv.cursor--
		}
	case "down", "j":
		if mv.cursor < len(m.Menu())-1 {
			mv.cursor++
		}
	case "enter", " ":
		if len(m.Menu()) == 0 {
			return nil
		}
		entry := m.Menu()[mv.cursor]
		switch entry.Type {
		case "content_type":
			return RequestOpenContentType(entry)
		case "static":
			return RequestOpenStaticPage(entry)
		}
	}
	return nil
}

func MenuWidth(totalWidth int) int {
	if totalWidth > layout.MaxLayoutWidth {
		return layout.MaxLayoutWidth
	}
	return totalWidth
}
