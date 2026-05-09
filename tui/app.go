package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jordi.codes/tui/components"
	"jordi.codes/tui/layout"
)

// ── tick messages ─────────────────────────────────────────────────────────────

type clockTickMsg time.Time

func doClockTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return clockTickMsg(t) })
}

// ── bubbletea interface ───────────────────────────────────────────────────────

func (ctx Context) Init() tea.Cmd {
	return tea.Batch(doClockTick(), fetchContribs(ctx.cfg.Site.GitHub))
}

func (ctx Context) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ctx.width = msg.Width
		ctx.height = msg.Height

	case clockTickMsg:
		ctx.now = time.Time(msg)
		return ctx, doClockTick()

	case contribsMsg:
		ctx.contribs = msg.counts
		return ctx, nil

	case components.OpenContentTypeMsg:
		ctx.pushView(components.NewListView(msg.Entry, ctx.cfg, ctx.listLayout))
		return ctx, nil

	case components.OpenStaticPageMsg:
		ctx.pushView(components.NewStaticDetailView(msg.Entry, ctx.width, ctx.height))
		return ctx, nil

	case components.OpenDetailMsg:
		ctx.pushView(components.NewDetailView(msg.Item.Title, msg.Item.Body, ctx.width, ctx.height))
		return ctx, nil

	case components.NavBackMsg:
		ctx.popView()
		return ctx, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ctx, tea.Quit
		}
	}

	if ctx.activeView != nil {
		cmd := ctx.activeView.Update(&ctx, msg)
		return ctx, cmd
	}

	return ctx, nil
}

func (ctx Context) View() string {
	if ctx.width < 40 || ctx.height < 10 {
		return "Terminal too small \u2014 please resize to at least 40\u00d710.\n"
	}

	if ctx.activeView == nil {
		return ""
	}

	body := ctx.activeView.Render(&ctx)
	footer := components.Footer{}.Render(&ctx)

	if _, ok := ctx.activeView.(*components.MenuView); ok {
		mw := components.MenuWidth(ctx.width)
		if mw < ctx.width {
			bodyHeight := ctx.height - layout.FooterHeight
			return layout.CenterBodyWithinScreen(ctx.width, ctx.height, bodyHeight, body, footer)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}
