package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/tui/layout"
)

type DetailView struct {
	title      string
	rawContent string
	viewport   viewport.Model
	err        error
}

func NewDetailView(title, body string, width, height int) *DetailView {
	rendered, err := RenderMarkdown(body, width)
	if err != nil {
		rendered = body
	}
	vp := viewport.New(width, layout.ViewportHeight(height))
	vp.SetContent(rendered)
	return &DetailView{
		title:      title,
		rawContent: body,
		viewport:   vp,
	}
}

func NewErrorDetailView(title string, err error) *DetailView {
	return &DetailView{
		title: title,
		err:   fmt.Errorf("could not read page: %w", err),
	}
}

func (d *DetailView) Render(m AppContext) string {
	if d.err != nil {
		m.SetHelpText("esc  back   q  quit")
		return ErrorView{ErrText: d.err.Error()}.Render(m)
	}
	m.SetHelpText("↑/↓  scroll   PgUp/PgDn   esc  back   q  quit")
	return DetailBody{
		Title:    d.title,
		Viewport: d.viewport,
	}.Render(m)
}

func (d *DetailView) Update(m AppContext, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.viewport.Width = msg.Width
		d.viewport.Height = layout.ViewportHeight(msg.Height)
		if d.rawContent != "" {
			if rendered, err := RenderMarkdown(d.rawContent, msg.Width); err == nil {
				d.viewport.SetContent(rendered)
			}
		}
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
