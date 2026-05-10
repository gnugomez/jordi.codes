package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/cms"
	"jordi.codes/router"
	"jordi.codes/tui/components"
)

type Context struct {
	site   *cms.Site
	width  int
	height int

	listLayout components.ListLayout

	now        time.Time
	remoteAddr string

	contribs map[string]int // date "2006-01-02" → count

	activeView components.ViewComponent
	viewStack  []components.ViewComponent

	helpText string

	initialPath string // URL path to navigate to on startup (empty = main menu)
}

type Option func(*Context)

func WithListLayout(l components.ListLayout) Option {
	return func(ctx *Context) {
		if l != nil {
			ctx.listLayout = l
		}
	}
}

// WithInitialPath sets a URL path that the TUI will navigate to immediately
// on startup. The path follows the same routing convention used by the HTTP
// server: "/about", "/projects", "/projects/jordi-codes".
// SSH sessions can use dot-notation usernames: "projects.jordi-codes".
func WithInitialPath(path string) Option {
	return func(ctx *Context) {
		ctx.initialPath = path
	}
}

func newContext(site *cms.Site, width, height int, remoteAddr string, opts ...Option) Context {
	ctx := Context{
		site:       site,
		width:      width,
		height:     height,
		now:        time.Now(),
		remoteAddr: remoteAddr,
		listLayout: components.StackedBoxListLayout{},
		activeView: components.NewMenuView(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&ctx)
		}
	}

	return ctx
}

func (ctx *Context) pushView(v components.ViewComponent) {
	ctx.viewStack = append(ctx.viewStack, ctx.activeView)
	ctx.activeView = v
}

func (ctx *Context) popView() {
	n := len(ctx.viewStack)
	if n > 0 {
		ctx.activeView = ctx.viewStack[n-1]
		ctx.viewStack = ctx.viewStack[:n-1]
	}
}

// initialNavCmd resolves ctx.initialPath to a tea.Cmd that pushes the
// matching view onto the stack. Returns nil if the path is empty or unknown.
func (ctx *Context) initialNavCmd() tea.Cmd {
	if ctx.initialPath == "" {
		return nil
	}
	r := router.New(ctx.site)
	route, ok := r.Resolve(ctx.initialPath)
	if !ok {
		return nil
	}
	switch route.Kind {
	case router.KindStatic:
		return components.RequestOpenStaticPage(route.Entry)
	case router.KindList:
		return components.RequestOpenContentType(route.Entry)
	case router.KindDetail:
		item, err := ctx.site.LoadContentItemBySlug(route.Entry, route.Slug)
		if err != nil {
			return nil
		}
		return components.RequestOpenDetail(item)
	}
	return nil
}

func (ctx *Context) SetHelpText(text string) { ctx.helpText = text }
func (ctx *Context) HelpText() string        { return ctx.helpText }

// ── AppContext accessors ───────────────────────────────────────────────────────

func (ctx *Context) Width() int                        { return ctx.width }
func (ctx *Context) Height() int                       { return ctx.height }
func (ctx *Context) Menu() []cms.MenuEntry             { return ctx.site.Menu }
func (ctx *Context) Subtitle() string                  { return ctx.site.Site.Subtitle }
func (ctx *Context) Now() time.Time                    { return ctx.now }
func (ctx *Context) Contribs() map[string]int          { return ctx.contribs }
func (ctx *Context) RemoteAddr() string                { return ctx.remoteAddr }
func (ctx *Context) ListLayout() components.ListLayout { return ctx.listLayout }
