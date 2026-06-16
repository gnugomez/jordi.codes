package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type MenuWidgetPanel struct {
	Width  int
	Height int
}

func (p MenuWidgetPanel) Render(m AppContext) string {
	t := m.Theme()
	width := p.Width
	height := p.Height
	const gap = 2
	innerWidth := width - gap

	greeting := widgetGreeting(t, innerWidth, m.Now())
	activity := widgetActivity(t, innerWidth, m.Contribs(), m.Now())
	visitor := widgetVisitor(t, innerWidth, m.RemoteAddr())

	stack := strings.Join([]string{greeting, activity, visitor}, "\n")
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, stack)
}

func widget(t *Theme, content string, outerWidth int) string {
	inner := outerWidth - 2
	if inner < 4 {
		inner = 4
	}
	return t.WidgetBorderStyle.Width(inner).Render(content)
}

func widgetLabel(t *Theme, s string) string {
	return t.WidgetTitleStyle.Render(strings.ToUpper(s))
}

var catPortrait = []string{
	`  ∧,,,∧`,
	` (• ⩊ •)`,
}

func widgetGreeting(t *Theme, outerWidth int, now time.Time) string {
	hour := now.Hour()
	var greeting string
	switch {
	case hour < 5:
		greeting = "up late?"
	case hour < 12:
		greeting = "good morning~"
	case hour < 17:
		greeting = "good afternoon~"
	case hour < 21:
		greeting = "good evening~"
	default:
		greeting = "still up?"
	}

	bubble := t.WidgetAccentStyle.Render("hey! "+greeting) + "\n" +
		t.SubtitleStyle.Render("my name is jordi 🫧") + "\n" +
		t.SubtitleStyle.Render("and my belief is that art can come from code")

	const catMinWidth = 32
	var row string
	if outerWidth >= catMinWidth {
		portrait := strings.Join(catPortrait, "\n")
		row = lipgloss.JoinHorizontal(lipgloss.Top,
			t.NewStyle().PaddingRight(2).Render(portrait),
			bubble,
		)
	} else {
		row = bubble
	}
	return widget(t, widgetLabel(t, "about me")+"\n\n"+row, outerWidth)
}

var contribColors = []lipgloss.Color{"#1C1917", "#92400E", "#C2410C", "#EA580C", "#F97316"}

func contribLevel(n int) int {
	switch {
	case n <= 0:
		return 0
	case n <= 3:
		return 1
	case n <= 6:
		return 2
	case n <= 9:
		return 3
	default:
		return 4
	}
}

func buildContribGrid(numWeeks int, gridStart time.Time, contribs map[string]int) [7][]int {
	var grid [7][]int
	for i := range grid {
		grid[i] = make([]int, numWeeks)
	}
	for w := 0; w < numWeeks; w++ {
		for d := 0; d < 7; d++ {
			day := gridStart.AddDate(0, 0, w*7+d)
			grid[int(day.Weekday())][w] = contribLevel(contribs[day.Format("2006-01-02")])
		}
	}
	return grid
}

func buildMonthHeader(t *Theme, numWeeks int, gridStart time.Time, labelWidth int) string {
	runes := make([]rune, numWeeks)
	for i := range runes {
		runes[i] = ' '
	}
	prevMonth := -1
	for w := 0; w < numWeeks; w++ {
		day := gridStart.AddDate(0, 0, w*7)
		if m := int(day.Month()); m != prevMonth {
			for j, c := range []rune(day.Format("Jan")) {
				if w+j < numWeeks {
					runes[w+j] = c
				}
			}
			prevMonth = m
		}
	}
	return strings.Repeat(" ", labelWidth) + t.MutedStyle.Render(string(runes))
}

func buildGridRows(t *Theme, grid [7][]int, numWeeks int) [7]strings.Builder {
	block := "█"
	var rows [7]strings.Builder
	for row := 0; row < 7; row++ {
		for w := 0; w < numWeeks; w++ {
			style := t.NewStyle().Foreground(contribColors[grid[row][w]])
			rows[row].WriteString(style.Render(block))
		}
	}
	return rows
}

func countActiveDays(contribs map[string]int) int {
	n := 0
	for _, v := range contribs {
		if v > 0 {
			n++
		}
	}
	return n
}

func activityContent(t *Theme, outerWidth int, contribs map[string]int, now time.Time) string {
	inner := outerWidth - 2
	if inner < 10 {
		inner = 10
	}

	const (
		labelWidth = 3
		cellWidth  = 1
	)
	numWeeks := (inner - labelWidth) / cellWidth
	if numWeeks > 53 {
		numWeeks = 53
	}

	daysToSat := (6 - int(now.Weekday()) + 7) % 7
	gridEnd := now.AddDate(0, 0, daysToSat)
	gridStart := gridEnd.AddDate(0, 0, -(numWeeks*7 - 1))

	grid := buildContribGrid(numWeeks, gridStart, contribs)
	monthHeader := buildMonthHeader(t, numWeeks, gridStart, labelWidth)
	rows := buildGridRows(t, grid, numWeeks)

	dayLabels := [7]string{"", "Mo", "", "We", "", "Fr", ""}
	var sb strings.Builder
	sb.WriteString(monthHeader + "\n")
	for row := 0; row < 7; row++ {
		if dayLabels[row] != "" {
			sb.WriteString(t.MutedStyle.Render(dayLabels[row]) + " ")
		} else {
			sb.WriteString("   ")
		}
		sb.WriteString(rows[row].String() + "\n")
	}

	var totalStr string
	if len(contribs) > 0 {
		totalStr = t.WidgetAccentStyle.Render(fmt.Sprintf("%d", countActiveDays(contribs))) + t.SubtitleStyle.Render(" active days in the last year")
	} else {
		totalStr = t.MutedStyle.Render("loading…")
	}

	return widgetLabel(t, "github activity") + "\n\n" +
		t.NewStyle().PaddingLeft((inner-labelWidth-numWeeks*cellWidth)/2).Render(strings.TrimRight(sb.String(), "\n")) + "\n\n" +
		t.NewStyle().Width(inner).Align(lipgloss.Center).Render(totalStr)
}

func widgetActivity(t *Theme, outerWidth int, contribs map[string]int, now time.Time) string {
	return widget(t, activityContent(t, outerWidth, contribs, now), outerWidth)
}

func widgetVisitor(t *Theme, outerWidth int, remoteAddr string) string {
	if remoteAddr == "" {
		remoteAddr = "unknown"
	}
	face := t.MutedStyle.Render("(✿◠‿◠)")
	line := face + "  " + t.SubtitleStyle.Render("connected from ") + t.WidgetAccentStyle.Render(remoteAddr)
	return widget(t, widgetLabel(t, "visitor")+"\n\n"+line, outerWidth)
}
