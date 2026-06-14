package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"jordi.codes/cms"
	"jordi.codes/tui/layout"
)

// ListLayout decouples list navigation and rendering from concrete view styles.
type ListLayout interface {
	VisibleItems(bodyHeight int) int
	Render(items []cms.ContentItem, cursor, offset, width, bodyHeight int) string
}

// StackedBoxListLayout renders content items as a vertical stack of bordered cards.
type StackedBoxListLayout struct{}

const (
	listCardGap         = 1
	listCardContentRows = 2
	listCardOuterRows   = listCardContentRows + 2
)

func (StackedBoxListLayout) VisibleItems(bodyHeight int) int {
	if bodyHeight <= 0 {
		return 1
	}
	rowsPerItem := listCardOuterRows + listCardGap
	if rowsPerItem <= 0 {
		return 1
	}
	if n := (bodyHeight + listCardGap) / rowsPerItem; n > 0 {
		return n
	}
	return 1
}

func (l StackedBoxListLayout) Render(items []cms.ContentItem, cursor, offset, width, bodyHeight int) string {
	if len(items) == 0 {
		return lipgloss.NewStyle().
			PaddingLeft(2).
			Width(width).
			Height(bodyHeight).
			Render("\n" + mutedStyle.Render("   No entries found."))
	}

	visible := l.VisibleItems(bodyHeight)
	end := offset + visible
	if end > len(items) {
		end = len(items)
	}

	cardWidth := width - 6
	if cardWidth < 16 {
		cardWidth = 16
	}

	normalCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMuted).
		Padding(0, 1).
		Foreground(colorNormal).
		Width(cardWidth)

	selectedCardStyle := normalCardStyle.
		BorderForeground(colorPrimary).
		Foreground(colorPrimary).
		Bold(true)

	previewStyle := lipgloss.NewStyle().Foreground(colorSubtitle)

	cards := make([]string, 0, end-offset)
	for i := offset; i < end; i++ {
		item := items[i]
		style := normalCardStyle
		if i == cursor {
			style = selectedCardStyle
		}

		title := ellipsize(oneLine(item.Title), cardWidth-4)
		preview := previewText(item)
		preview = previewStyle.Render(ellipsize(oneLine(preview), cardWidth-4))

		cards = append(cards, style.Render(title+"\n"+preview))
	}

	stack := strings.Join(cards, "\n")
	return lipgloss.NewStyle().
		PaddingLeft(2).
		Width(width).
		Height(bodyHeight).
		Render(stack)
}

// ListBody renders a section header followed by the list layout body.
type ListBody struct {
	Title  string
	Items  []cms.ContentItem
	Cursor int
	Offset int
	Layout ListLayout
}

func (lb ListBody) Render(m AppContext) string {
	head := SectionHeader{Title: lb.Title}.Render(m)
	l := lb.Layout
	if l == nil {
		l = StackedBoxListLayout{}
	}
	bodyHeight := layout.BodyHeight(m.Height())
	body := l.Render(lb.Items, lb.Cursor, lb.Offset, m.Width(), bodyHeight)
	return lipgloss.JoinVertical(lipgloss.Left, head, body)
}

func previewText(item cms.ContentItem) string {
	if strings.TrimSpace(item.Excerpt) != "" {
		return item.Excerpt
	}
	body := strings.TrimSpace(item.Body)
	if body == "" {
		return "No preview available."
	}
	return body
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func ellipsize(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return string(r[:max-1]) + "…"
}
