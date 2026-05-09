package layout

import "github.com/charmbracelet/lipgloss"

const (
	MaxLayoutWidth = 120
	MinSplitWidth  = 90
	LeftPanelWidth = 42
	HeaderHeight   = 2 // section title + separator
	FooterHeight   = 2 // separator + help/clock line
)

func CenterBodyWithinScreen(totalWidth, totalHeight, bodyHeight int, body, footer string) string {
	centered := lipgloss.Place(totalWidth, bodyHeight, lipgloss.Center, lipgloss.Top, body)
	return lipgloss.JoinVertical(lipgloss.Left, centered, footer)
}

func PlaceCentered(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// ViewportHeight returns the viewport height given the total terminal height.
func ViewportHeight(totalHeight int) int {
	if h := totalHeight - HeaderHeight - FooterHeight; h > 0 {
		return h
	}
	return 1
}
