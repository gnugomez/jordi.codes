package components

import (
	"fmt"
	"io/fs"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/cms"
	"jordi.codes/tui/layout"
)

// DetailView is a scrollable view of a single content item's markdown body.
type DetailView struct {
	item     cms.ContentItem
	viewport viewport.Model
	content  MarkdownContent
	err      error
}

// NewDetailView creates a detail view for a content item.
func NewDetailView(item cms.ContentItem, width, height int, fsys fs.FS) *DetailView {
	mc := MarkdownContent{Body: item.Body, FS: fsys, ContentDir: item.ContentDir}
	vp := mc.Viewport(width, layout.ViewportHeight(height))
	return &DetailView{item: item, viewport: vp, content: mc}
}

// NewErrorDetailView creates a detail view that displays an error.
func NewErrorDetailView(title string, err error) *DetailView {
	return &DetailView{
		item: cms.ContentItem{Title: title},
		err:  fmt.Errorf("could not read page: %w", err),
	}
}

func (d *DetailView) Render(m AppContext) string {
	if d.err != nil {
		m.SetHelpText("esc  back   q  quit")
		return ErrorView{ErrText: d.err.Error()}.Render(m)
	}
	m.SetHelpText("↑/↓  scroll   PgUp/PgDn   esc  back   q  quit")
	return ContentBody{
		Title:    d.item.Title,
		Viewport: d.viewport,
	}.Render(m)
}

func (d *DetailView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.viewport.Width = msg.Width
		d.viewport.Height = layout.ViewportHeight(msg.Height)
		d.content.RefreshViewport(&d.viewport, msg.Width)
		return nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return tea.Quit
		case "esc":
			return RequestNavBack()
		}
	}
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return cmd
}
