package components

import (
	tea "github.com/charmbracelet/bubbletea"

	"jordi.codes/cms"
)

// Navigation messages — view components emit these to request transitions.
// The app layer creates the target component and pushes it onto the view stack.

type OpenContentTypeMsg struct{ Entry cms.MenuEntry }
type OpenStaticPageMsg struct{ Entry cms.MenuEntry }
type OpenDetailMsg struct{ Item cms.ContentItem }
type NavBackMsg struct{}

func RequestOpenContentType(entry cms.MenuEntry) tea.Cmd {
	return func() tea.Msg { return OpenContentTypeMsg{Entry: entry} }
}

func RequestOpenStaticPage(entry cms.MenuEntry) tea.Cmd {
	return func() tea.Msg { return OpenStaticPageMsg{Entry: entry} }
}

func RequestOpenDetail(item cms.ContentItem) tea.Cmd {
	return func() tea.Msg { return OpenDetailMsg{Item: item} }
}

func RequestNavBack() tea.Cmd {
	return func() tea.Msg { return NavBackMsg{} }
}
