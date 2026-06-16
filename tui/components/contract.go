package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/cms"
)

// AppContext is the minimal interface view components need from the app context.
type AppContext interface {
	Width() int
	Height() int
	SetHelpText(string)
	HelpText() string
	Menu() []cms.MenuEntry
	Subtitle() string
	Now() time.Time
	Contribs() map[string]int
	RemoteAddr() string
	ListLayout() ListLayout
	Theme() *Theme
}

// Renderer is the contract for any component that can draw itself.
type Renderer interface {
	Render(AppContext) string
}

// ViewComponent is the contract for a top-level view.
// Each component owns its own state and handles all message types.
type ViewComponent interface {
	Renderer
	Update(m AppContext, msg tea.Msg) tea.Cmd
}
