package components

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
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
	loading bool
	spinner spinner.Model
	err     error
	preview viewport.Model
	fsys    fs.FS
}

// previewImagesReadyMsg is delivered when background image rendering for the
// preview pane completes. idx identifies the item the render belongs to.
type previewImagesReadyMsg struct {
	idx     int
	content string
}

// contentLoadedMsg is delivered when async content loading completes.
type contentLoadedMsg struct {
	items []cms.ContentItem
	err   error
}

func NewListView(entry cms.MenuEntry, site *cms.Site, lo ListLayout, width, height int) (*ListView, tea.Cmd) {
	ct := site.ContentTypeByName(entry.ContentType)
	isAsync := ct != nil && strings.EqualFold(strings.TrimSpace(ct.Source), "github_pinned")

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)
	lv := &ListView{title: entry.Label, layout: lo, fsys: site.FS(), spinner: s}

	if isAsync {
		lv.loading = true
		loadCmd := func() tea.Msg {
			items, err := site.LoadMenuContentItems(entry)
			if err != nil {
				return contentLoadedMsg{err: fmt.Errorf("could not load %s: %w", entry.Label, err)}
			}
			return contentLoadedMsg{items: items}
		}
		return lv, tea.Batch(loadCmd, lv.spinner.Tick)
	}

	items, err := site.LoadMenuContentItems(entry)
	if err != nil {
		lv.err = fmt.Errorf("could not load %s: %w", entry.Label, err)
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

func (lv *ListView) previewWidth(totalWidth int) int {
	listW := lv.listWidth(totalWidth)
	extra := 1 // vertical separator between list and preview
	previewW := totalWidth - listW - extra
	if previewW < 20 {
		previewW = 20
	}
	return previewW
}

func (lv *ListView) updatePreview(totalWidth, totalHeight int) tea.Cmd {
	previewW := lv.previewWidth(totalWidth)
	vpH := layout.ViewportHeight(totalHeight)

	if len(lv.items) == 0 || lv.cursor < 0 || lv.cursor >= len(lv.items) {
		lv.preview = viewport.New(previewW, vpH)
		return nil
	}

	mc := MarkdownContent{Body: lv.items[lv.cursor].Body, FS: lv.fsys, ContentDir: lv.items[lv.cursor].ContentDir}
	lv.preview = mc.ViewportFast(previewW, vpH)

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

func (lv *ListView) renderLoading(m AppContext) string {
	bodyHeight := m.Height() - layout.HeaderHeight - layout.FooterHeight

	loadingText := lv.spinner.View() + mutedStyle.Render(" Loading…")

	if !lv.wideEnough(m.Width()) {
		head := SectionHeader{Title: lv.title}.Render(m)
		body := lipgloss.NewStyle().
			Width(m.Width()).
			Height(bodyHeight).
			Align(lipgloss.Center, lipgloss.Center).
			Render(loadingText)
		return lipgloss.JoinVertical(lipgloss.Left, head, body)
	}

	listW := lv.listWidth(m.Width())
	previewW := lv.previewWidth(m.Width())
	head := lv.renderWideHeader(m, listW)

	listPanel := lipgloss.NewStyle().
		Width(listW).
		MaxWidth(listW).
		Height(bodyHeight).
		MaxHeight(bodyHeight).
		Render("")
	previewPanel := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(previewW).
		MaxWidth(previewW).
		Height(bodyHeight).
		MaxHeight(bodyHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(loadingText)

	panels := []string{listPanel, lv.verticalRule(bodyHeight), previewPanel}

	body := lipgloss.JoinHorizontal(lipgloss.Top, panels...)
	return lipgloss.JoinVertical(lipgloss.Left, head, body)
}

func (lv *ListView) renderWideHeader(m AppContext, listW int) string {
	titleLine := lipgloss.NewStyle().
		PaddingLeft(2).
		Render(mutedStyle.Render("◈  ") + headerStyle.Render(lv.title))

	separator := []rune(strings.Repeat("─", m.Width()))

	firstDivider := listW
	if firstDivider >= 0 && firstDivider < len(separator) {
		separator[firstDivider] = '┬'
	}

	return titleLine + "\n" + mutedStyle.Render(string(separator))
}

func (lv *ListView) verticalRule(height int) string {
	if height <= 0 {
		return ""
	}
	lines := make([]string, height)
	for i := range lines {
		lines[i] = "│"
	}
	return mutedStyle.Render(strings.Join(lines, "\n"))
}

func (lv *ListView) Render(m AppContext) string {
	if lv.loading {
		m.SetHelpText("loading…")
		return lv.renderLoading(m)
	}
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
	previewW := lv.previewWidth(m.Width())
	bodyHeight := m.Height() - layout.HeaderHeight - layout.FooterHeight

	head := lv.renderWideHeader(m, listW)

	l := lv.layout
	if l == nil {
		l = StackedBoxListLayout{}
	}
	listPanel := lipgloss.NewStyle().
		Width(listW).
		MaxWidth(listW).
		Height(bodyHeight).
		MaxHeight(bodyHeight).
		Render(l.Render(lv.items, lv.cursor, lv.offset, listW, bodyHeight))
	previewPanel := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(previewW).
		MaxWidth(previewW).
		Height(bodyHeight).
		MaxHeight(bodyHeight).
		Render(lv.preview.View())
	panels := []string{listPanel, lv.verticalRule(bodyHeight), previewPanel}

	body := lipgloss.JoinHorizontal(lipgloss.Top, panels...)
	return lipgloss.JoinVertical(lipgloss.Left, head, body)
}

func (lv *ListView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(contentLoadedMsg); ok {
		lv.loading = false
		if msg.err != nil {
			lv.err = msg.err
		} else {
			lv.items = msg.items
		}
		return lv.updatePreview(m.Width(), m.Height())
	}

	if lv.loading {
		var cmd tea.Cmd
		lv.spinner, cmd = lv.spinner.Update(msg)
		return cmd
	}

	if lv.wideEnough(m.Width()) {
		if cmd, handled := lv.handleWideMessage(m, msg); handled {
			return cmd
		}
	}

	cmd, rerenderPreview := lv.handleNavigationMessage(m, msg)
	if cmd != nil {
		return cmd
	}
	if rerenderPreview && lv.wideEnough(m.Width()) {
		return lv.updatePreview(m.Width(), m.Height())
	}
	return nil
}

func (lv *ListView) handleWideMessage(m AppContext, msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case contentLoadedMsg:
		return nil, false // handled in Update
	case tea.WindowSizeMsg:
		return lv.updatePreview(m.Width(), m.Height()), true
	case previewImagesReadyMsg:
		if msg.idx == lv.cursor {
			lv.preview.SetContent(msg.content)
		}
		return nil, true
	case tea.KeyMsg:
		// Let navigation keys (up/down/esc/enter) be handled by handleNavigationMessage.
		return nil, false
	default:
		var cmd tea.Cmd
		lv.preview, cmd = lv.preview.Update(msg)
		return cmd, true
	}
}

func (lv *ListView) handleNavigationMessage(m AppContext, msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}

	moved := false
	switch keyMsg.String() {
	case "q":
		return tea.Quit, false
	case "esc":
		return RequestNavBack(), false
	case "up", "k":
		if lv.cursor > 0 {
			lv.cursor--
			moved = true
			if lv.cursor < lv.offset {
				lv.offset = lv.cursor
			}
		}
	case "down", "j":
		if lv.cursor < len(lv.items)-1 {
			lv.cursor++
			moved = true
			if vis := lv.visibleItems(m.Height()); lv.cursor >= lv.offset+vis {
				lv.offset = lv.cursor - vis + 1
			}
		}
	case "enter", " ":
		if len(lv.items) > 0 {
			return RequestOpenDetail(lv.items[lv.cursor]), false
		}
	}

	return nil, moved
}

func (lv *ListView) visibleItems(height int) int {
	bodyHeight := layout.BodyHeight(height)
	if lv.layout == nil {
		return 1
	}
	return lv.layout.VisibleItems(bodyHeight)
}
