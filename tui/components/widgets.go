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
	width := p.Width
	height := p.Height
	const gap = 2
	innerWidth := width - gap

	greeting := widgetGreeting(innerWidth, m.Now())

	const minSideBySide = 70
	var actBlock, calBlock string
	if innerWidth >= minSideBySide {
		calW := 28
		actW := innerWidth - calW
		calC := calendarContent(calW, m.Now())
		actC := activityContent(actW, m.Contribs(), m.Now())
		calH := strings.Count(calC, "\n") + 1
		actH := strings.Count(actC, "\n") + 1
		maxH := calH
		if actH > maxH {
			maxH = actH
		}
		midRow := lipgloss.JoinHorizontal(lipgloss.Top,
			widgetH(calC, calW, maxH),
			widgetH(actC, actW, maxH),
		)
		actBlock = midRow
		calBlock = ""
	} else {
		actBlock = widgetActivity(innerWidth, m.Contribs(), m.Now())
		calBlock = widgetCalendar(innerWidth, m.Now())
	}

	visitor := widgetVisitor(innerWidth, m.RemoteAddr())

	blockH := func(s string) int { return strings.Count(s, "\n") + 1 }
	const sep = 1
	type block struct{ s string }
	var all []block
	all = append(all, block{greeting})
	all = append(all, block{actBlock})
	if calBlock != "" {
		all = append(all, block{calBlock})
	}
	all = append(all, block{visitor})

	for n := len(all); n >= 1; n-- {
		total := 0
		for i := 0; i < n; i++ {
			total += blockH(all[i].s)
			if i < n-1 {
				total += sep
			}
		}
		if total <= height {
			parts := make([]string, n)
			for i := 0; i < n; i++ {
				parts[i] = all[i].s
			}
			stack := strings.Join(parts, "\n")
			return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, stack)
		}
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, "")
}

func widget(content string, outerWidth int) string {
	inner := outerWidth - 2
	if inner < 4 {
		inner = 4
	}
	return widgetBorderStyle.Width(inner).Render(content)
}

func widgetH(content string, outerWidth, innerHeight int) string {
	inner := outerWidth - 2
	if inner < 4 {
		inner = 4
	}
	return widgetBorderStyle.Width(inner).Height(innerHeight).Render(content)
}

func widgetLabel(s string) string {
	return widgetTitleStyle.Render(strings.ToUpper(s))
}

var catPortrait = []string{
	` ∧,,,∧`,
	` (• ⩊ •)`,
}

func widgetGreeting(outerWidth int, now time.Time) string {
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

	bubble := widgetAccentStyle.Render("hey! "+greeting) + "\n" +
		subtitleStyle.Render("i'm jordi  —  software engineer,") + "\n" +
		subtitleStyle.Render("coffee addict & terminal dweller.")

	const catMinWidth = 32
	var row string
	if outerWidth >= catMinWidth {
		portrait := strings.Join(catPortrait, "\n")
		row = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().PaddingRight(2).Render(portrait),
			bubble,
		)
	} else {
		row = bubble
	}
	return widget(widgetLabel("about me")+"\n\n"+row, outerWidth)
}

func calendarContent(outerWidth int, now time.Time) string {
	month := strings.ToUpper(now.Format("Jan"))
	year := now.Format("2006")
	day := now.Format("2")
	weekday := now.Format("Monday")

	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	startOffset := int(firstDay.Weekday())

	header := mutedStyle.Render("Su Mo Tu We Th Fr Sa")
	var weeks []string
	col := 0
	var week strings.Builder
	for i := 0; i < startOffset; i++ {
		week.WriteString("   ")
		col++
	}
	for d := 1; d <= daysInMonth; d++ {
		if col == 7 {
			weeks = append(weeks, week.String())
			week.Reset()
			col = 0
		}
		cell := fmt.Sprintf("%2d", d)
		if d == now.Day() {
			week.WriteString(widgetAccentStyle.Render(cell) + " ")
		} else {
			week.WriteString(normalItemStyle.Render(cell) + " ")
		}
		col++
	}
	if week.Len() > 0 {
		weeks = append(weeks, week.String())
	}

	headline := widgetAccentStyle.Render(day) +
		mutedStyle.Render("  "+weekday) + "\n" +
		subtitleStyle.Render(month+" "+year)

	inner := outerWidth - 2
	gridBlock := header + "\n" + strings.Join(weeks, "\n")
	centeredGrid := lipgloss.NewStyle().Width(inner).Align(lipgloss.Center).Render(gridBlock)
	centeredHeadline := lipgloss.NewStyle().Width(inner).Align(lipgloss.Center).Render(headline)

	content := widgetLabel("calendar") + "\n\n" + centeredHeadline + "\n\n" + centeredGrid
	return content
}

func widgetCalendar(outerWidth int, now time.Time) string {
	return widget(calendarContent(outerWidth, now), outerWidth)
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

func activityContent(outerWidth int, contribs map[string]int, now time.Time) string {
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

	type cell struct{ level int }
	grid := make([][]cell, 7)
	for i := range grid {
		grid[i] = make([]cell, numWeeks)
	}

	todayWeekday := int(now.Weekday())
	daysToSat := (6 - todayWeekday + 7) % 7
	gridEnd := now.AddDate(0, 0, daysToSat)
	gridStart := gridEnd.AddDate(0, 0, -(numWeeks*7 - 1))

	for w := 0; w < numWeeks; w++ {
		for d := 0; d < 7; d++ {
			offset := w*7 + d
			day := gridStart.AddDate(0, 0, offset)
			row := int(day.Weekday())
			key := day.Format("2006-01-02")
			n := contribs[key]
			grid[row][w] = cell{level: contribLevel(n)}
		}
	}

	block := "█"
	var rows [7]strings.Builder
	for row := 0; row < 7; row++ {
		for w := 0; w < numWeeks; w++ {
			lvl := grid[row][w].level
			style := lipgloss.NewStyle().Foreground(contribColors[lvl])
			rows[row].WriteString(style.Render(block))
		}
	}

	dayLabels := [7]string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	var sb strings.Builder
	for row := 0; row < 7; row++ {
		var label string
		if row%2 == 0 {
			label = mutedStyle.Render(dayLabels[row]) + " "
		} else {
			label = "   "
		}
		sb.WriteString(label + rows[row].String() + "\n")
	}

	activeDays := 0
	for _, v := range contribs {
		if v > 0 {
			activeDays++
		}
	}
	var totalStr string
	if len(contribs) > 0 {
		totalStr = widgetAccentStyle.Render(fmt.Sprintf("%d", activeDays)) + subtitleStyle.Render(" active days in the last year")
	} else {
		totalStr = mutedStyle.Render("loading…")
	}

	return widgetLabel("github activity") + "\n\n" +
		lipgloss.NewStyle().PaddingLeft((inner-labelWidth-numWeeks*cellWidth)/2).Render(strings.TrimRight(sb.String(), "\n")) + "\n\n" +
		lipgloss.NewStyle().Width(inner).Align(lipgloss.Center).Render(totalStr)
}

func widgetActivity(outerWidth int, contribs map[string]int, now time.Time) string {
	return widget(activityContent(outerWidth, contribs, now), outerWidth)
}

func widgetVisitor(outerWidth int, remoteAddr string) string {
	if remoteAddr == "" {
		remoteAddr = "unknown"
	}
	face := mutedStyle.Render("(✿◠‿◠)")
	line := face + "  " + subtitleStyle.Render("connected from ") + widgetAccentStyle.Render(remoteAddr)
	return widget(widgetLabel("visitor")+"\n\n"+line, outerWidth)
}
