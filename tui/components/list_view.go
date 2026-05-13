package components

import (
	"fmt"
	"io/fs"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jordi.codes/cms"
	"jordi.codes/tui/layout"
)

const minSplitWidth = 100

type ListView struct {
	title   string
	items   []cms.ContentItem
	cursor  int
	offset  int
	layout  ListLayout
	err     error
	preview viewport.Model
	prevIdx int // cursor index last rendered in preview
	prevW   int // width used for last preview render
	fsys    fs.FS
}

// previewImagesReadyMsg is delivered when background image rendering for the
// preview pane completes. idx identifies the item the render belongs to.
type previewImagesReadyMsg struct {
	idx     int
	content string
}

func NewListView(entry cms.MenuEntry, site *cms.Site, lo ListLayout, width, height int) (*ListView, tea.Cmd) {
	lv := &ListView{title: entry.Label, layout: lo, prevIdx: -1, fsys: site.FS()}

	ct := cms.FindContentType(site, entry.ContentType)
	if ct == nil {
		lv.err = fmt.Errorf("content type %q not found in config", entry.ContentType)
		return lv, nil
	}

	items, err := site.LoadContentItems(*ct)
	if err != nil {
		lv.err = fmt.Errorf("could not load %s: %w", ct.DisplayName, err)
		return lv, nil
	}

	lv.items = items
	return lv, lv.updatePreview(width, height)
}

func (lv *ListView) wideEnough(w int) bool {
	return w >= minSplitWidth
}

func (lv *ListView) listWidth(totalWidth int) int {
	w := totalWidth / 4
	if w < 30 {
		w = 30
	}
	return w
}

func (lv *ListView) updatePreview(totalWidth, totalHeight int) tea.Cmd {
	if len(lv.items) == 0 {
		return nil
	}
	listW := lv.listWidth(totalWidth)
	previewW := totalWidth - listW - 3 // 3 for border + padding
	vpH := layout.ViewportHeight(totalHeight)

	if lv.cursor == lv.prevIdx && previewW == lv.prevW {
		return nil
	}

	mc := MarkdownContent{Body: lv.items[lv.cursor].Body, FS: lv.fsys, ContentDir: lv.items[lv.cursor].ContentDir}
	lv.preview = mc.ViewportFast(previewW, vpH)
	lv.prevIdx = lv.cursor
	lv.prevW = previewW

	idx := lv.cursor
	if imgCmd := mc.RenderImagesCmd(previewW); imgCmd != nil {
		return func() tea.Msg {
			if m, ok := imgCmd().(imagesReadyMsg); ok {
				return previewImagesReadyMsg{idx: idx, content: m.content}
			}
			return nil
		}
	}
	return nil
}

func (lv *ListView) Render(m AppContext) string {
	if lv.err != nil {
		m.SetHelpText("esc  back   q  quit")
		return ErrorView{ErrText: lv.err.Error()}.Render(m)
	}
	m.SetHelpText("↑/↓  navigate   enter  open   esc  back   q  quit")

	if !lv.wideEnough(m.Width()) {
		return ListBody{
			Title:  lv.title,
			Items:  lv.items,
			Cursor: lv.cursor,
			Offset: lv.offset,
			Layout: lv.layout,
		}.Render(m)
	}

	// Split layout: list on left, preview on right
	listW := lv.listWidth(m.Width())
	bodyHeight := m.Height() - layout.HeaderHeight - layout.FooterHeight

	head := SectionHeader{Title: lv.title}.Render(m)

	l := lv.layout
	if l == nil {
		l = StackedBoxListLayout{}
	}
	listPanel := l.Render(lv.items, lv.cursor, lv.offset, listW, bodyHeight)

	lv.updatePreview(m.Width(), m.Height())

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(colorMuted).
		PaddingLeft(1).
		Height(bodyHeight)

	previewPanel := borderStyle.Render(lv.preview.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, previewPanel)
	return lipgloss.JoinVertical(lipgloss.Left, head, body)
}

func (lv *ListView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	return lv.update(m, msg)
}

func (lv *ListView) update(m AppContext, msg tea.Msg) tea.Cmd {
	if lv.wideEnough(m.Width()) {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			lv.prevIdx = -1 // force re-render
			return lv.updatePreview(m.Width(), m.Height())
		case previewImagesReadyMsg:
			if msg.idx == lv.cursor {
				lv.preview.SetContent(msg.content)
			}
			return nil
		case tea.KeyMsg:
			// handle navigation below
			_ = msg
		default:
			var cmd tea.Cmd
			lv.preview, cmd = lv.preview.Update(msg)
			return cmd
		}
	}

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
	if lv.wideEnough(m.Width()) {
		return lv.updatePreview(m.Width(), m.Height())
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
