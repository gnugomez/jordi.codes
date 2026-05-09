package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/cms"
	"jordi.codes/tui/layout"
)

type ListView struct {
	title  string
	items  []cms.ContentItem
	cursor int
	offset int
	layout ListLayout
	err    error
}

func NewListView(entry cms.MenuEntry, cfg *cms.Config, lo ListLayout) *ListView {
	lv := &ListView{title: entry.Label, layout: lo}

	ct := cms.FindContentType(cfg, entry.ContentType)
	if ct == nil {
		lv.err = fmt.Errorf("content type %q not found in config", entry.ContentType)
		return lv
	}

	items, err := cms.LoadContentItems(*ct)
	if err != nil {
		lv.err = fmt.Errorf("could not load %s: %w", ct.DisplayName, err)
		return lv
	}

	lv.items = items
	return lv
}

func (lv *ListView) Render(m AppContext) string {
	if lv.err != nil {
		m.SetHelpText("esc  back   q  quit")
		return ErrorView{ErrText: lv.err.Error()}.Render(m)
	}
	m.SetHelpText("↑/↓  navigate   enter  open   esc  back   q  quit")
	return ListBody{
		Title:  lv.title,
		Items:  lv.items,
		Cursor: lv.cursor,
		Offset: lv.offset,
		Layout: lv.layout,
	}.Render(m)
}

func (lv *ListView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch keyMsg.String() {
	case "q":
		return tea.Quit
	case "esc":
		return RequestNavBack()
	case "up", "k":
		if lv.cursor > 0 {
			lv.cursor--
			if lv.cursor < lv.offset {
				lv.offset = lv.cursor
			}
		}
	case "down", "j":
		if lv.cursor < len(lv.items)-1 {
			lv.cursor++
			if vis := lv.visibleItems(m.Height()); lv.cursor >= lv.offset+vis {
				lv.offset = lv.cursor - vis + 1
			}
		}
	case "enter", " ":
		if len(lv.items) > 0 {
			return RequestOpenDetail(lv.items[lv.cursor])
		}
	}
	return nil
}

func (lv *ListView) visibleItems(height int) int {
	bodyHeight := height - layout.HeaderHeight - layout.FooterHeight
	if lv.layout == nil {
		return 1
	}
	return lv.layout.VisibleItems(bodyHeight)
}
