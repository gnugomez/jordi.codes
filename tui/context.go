package tui

import (
	"time"

	"jordi.codes/cms"
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
}

type Option func(*Context)

func WithListLayout(l components.ListLayout) Option {
	return func(ctx *Context) {
		if l != nil {
			ctx.listLayout = l
		}
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
