package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jordi.codes/cms"
	"jordi.codes/tui/components"
	"jordi.codes/tui/layout"
)

// ── tick messages ─────────────────────────────────────────────────────────────

type clockTickMsg time.Time

func doClockTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return clockTickMsg(t) })
}

// ── bubbletea model ───────────────────────────────────────────────────────────

// Model is the bubbletea tea.Model. It owns a Context which implements
// AppContext and is passed to all components for rendering and updates.
type Model struct {
	ctx Context
}

// NewModel creates the application model for a new SSH session.
func NewModel(site *cms.Site, width, height int, remoteAddr string, opts ...Option) Model {
	return Model{ctx: newContext(site, width, height, remoteAddr, opts...)}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(doClockTick(), fetchContribs(m.ctx.site.Site.GitHub))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.width = msg.Width
		m.ctx.height = msg.Height

	case clockTickMsg:
		m.ctx.now = time.Time(msg)
		return m, doClockTick()

	case contribsMsg:
		m.ctx.contribs = msg.counts
		return m, nil

	case components.OpenContentTypeMsg:
		m.ctx.pushView(components.NewListView(msg.Entry, m.ctx.site, m.ctx.listLayout))
		return m, nil

	case components.OpenStaticPageMsg:
		raw, err := m.ctx.site.LoadStaticPage(msg.Entry.Path)
		if err != nil {
			m.ctx.pushView(components.NewErrorDetailView(msg.Entry.Label, err))
		} else {
			m.ctx.pushView(components.NewDetailView(msg.Entry.Label, raw, m.ctx.width, m.ctx.height))
		}
		return m, nil

	case components.OpenDetailMsg:
		m.ctx.pushView(components.NewDetailView(msg.Item.Title, msg.Item.Body, m.ctx.width, m.ctx.height))
		return m, nil

	case components.NavBackMsg:
		m.ctx.popView()
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.ctx.activeView != nil {
		cmd := m.ctx.activeView.Update(&m.ctx, msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.ctx.width < 40 || m.ctx.height < 10 {
		return "Terminal too small \u2014 please resize to at least 40\u00d710.\n"
	}

	if m.ctx.activeView == nil {
		return ""
	}

	body := m.ctx.activeView.Render(&m.ctx)
	footer := components.Footer{}.Render(&m.ctx)

	if _, ok := m.ctx.activeView.(*components.MenuView); ok {
		mw := components.MenuWidth(m.ctx.width)
		if mw < m.ctx.width {
			bodyHeight := m.ctx.height - layout.FooterHeight
			return layout.CenterBodyWithinScreen(m.ctx.width, m.ctx.height, bodyHeight, body, footer)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}
