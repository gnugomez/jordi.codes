package components

import (
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
)

const (
	hexPrimary  = "#F97316"
	hexAccent   = "#FB923C"
	hexMuted    = "#78716C"
	hexNormal   = "#E7E5E4"
	hexClock    = "#FDBA74"
	hexSubtitle = "#A8A29E"
	hexError    = "#EF4444"
	hexBright   = "#FDE68A"
	hexAmber    = "#FBBF24"
	hexDarkBg   = "#1C1917"
	hexCodeBg   = "#292524"
)

var (
	colorPrimary  = lipgloss.Color(hexPrimary)
	colorAccent   = lipgloss.Color(hexAccent)
	colorMuted    = lipgloss.Color(hexMuted)
	colorNormal   = lipgloss.Color(hexNormal)
	colorClock    = lipgloss.Color(hexClock)
	colorSubtitle = lipgloss.Color(hexSubtitle)
	colorError    = lipgloss.Color(hexError)
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorSubtitle)

	headerStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(colorNormal)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	clockStyle = lipgloss.NewStyle().
			Foreground(colorClock)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	asciiGlowStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	widgetBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted)

	widgetTitleStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Bold(false)

	widgetAccentStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)
)

func sp(s string) *string { return &s }
func bp(b bool) *bool     { return &b }
func up(u uint) *uint     { return &u }

var amberStyle = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       sp(hexNormal),
		},
		Margin: up(2),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:  sp(hexClock),
			Italic: bp(true),
		},
		Indent:      up(1),
		IndentToken: sp("│ "),
	},
	List: ansi.StyleList{LevelIndent: 2},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       sp(hexPrimary),
			Bold:        bp(true),
		},
	},
	H1:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: " ", Suffix: " ", Color: sp(hexDarkBg), BackgroundColor: sp(hexPrimary), Bold: bp(true), Upper: bp(true)}},
	H2:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "## ", Color: sp(hexAccent), Bold: bp(true)}},
	H3:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "### ", Color: sp(hexClock), Bold: bp(true)}},
	H4:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "#### ", Color: sp(hexClock), Bold: bp(true)}},
	H5:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "##### ", Color: sp(hexClock), Bold: bp(true)}},
	H6:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "###### ", Color: sp(hexClock), Bold: bp(true)}},
	Text:           ansi.StylePrimitive{},
	Strikethrough:  ansi.StylePrimitive{CrossedOut: bp(true)},
	Emph:           ansi.StylePrimitive{Italic: bp(true), Color: sp(hexClock)},
	Strong:         ansi.StylePrimitive{Bold: bp(true), Color: sp(hexAccent)},
	HorizontalRule: ansi.StylePrimitive{Color: sp(hexMuted), Format: "\n--------\n"},
	Item:           ansi.StylePrimitive{BlockPrefix: "• "},
	Enumeration:    ansi.StylePrimitive{BlockPrefix: ". "},
	Task:           ansi.StyleTask{Ticked: "[✓] ", Unticked: "[ ] "},
	Link:           ansi.StylePrimitive{Color: sp(hexClock), Underline: bp(true)},
	LinkText:       ansi.StylePrimitive{Color: sp(hexAccent), Bold: bp(true)},
	Image:          ansi.StylePrimitive{Color: sp(hexClock), Underline: bp(true)},
	ImageText:      ansi.StylePrimitive{Color: sp(hexAccent), Format: "Image: {{.text}} →"},
	Code:           ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: " ", Suffix: " ", Color: sp(hexAmber), BackgroundColor: sp(hexCodeBg)}},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Color: sp(hexNormal)}, Margin: up(2)},
		Chroma: &ansi.Chroma{
			Text:                ansi.StylePrimitive{Color: sp(hexNormal)},
			Comment:             ansi.StylePrimitive{Color: sp(hexMuted)},
			CommentPreproc:      ansi.StylePrimitive{Color: sp(hexClock)},
			Keyword:             ansi.StylePrimitive{Color: sp(hexPrimary)},
			KeywordReserved:     ansi.StylePrimitive{Color: sp(hexPrimary)},
			KeywordNamespace:    ansi.StylePrimitive{Color: sp(hexAccent)},
			KeywordType:         ansi.StylePrimitive{Color: sp(hexClock)},
			Operator:            ansi.StylePrimitive{Color: sp(hexSubtitle)},
			Punctuation:         ansi.StylePrimitive{Color: sp(hexSubtitle)},
			Name:                ansi.StylePrimitive{Color: sp(hexNormal)},
			NameBuiltin:         ansi.StylePrimitive{Color: sp(hexAmber)},
			NameTag:             ansi.StylePrimitive{Color: sp(hexPrimary)},
			NameAttribute:       ansi.StylePrimitive{Color: sp(hexClock)},
			NameClass:           ansi.StylePrimitive{Color: sp(hexAmber), Underline: bp(true), Bold: bp(true)},
			NameConstant:        ansi.StylePrimitive{Color: sp(hexAmber)},
			NameDecorator:       ansi.StylePrimitive{Color: sp(hexAccent)},
			NameFunction:        ansi.StylePrimitive{Color: sp(hexAmber)},
			LiteralNumber:       ansi.StylePrimitive{Color: sp(hexClock)},
			LiteralString:       ansi.StylePrimitive{Color: sp(hexBright)},
			LiteralStringEscape: ansi.StylePrimitive{Color: sp("#F59E0B")},
			GenericDeleted:      ansi.StylePrimitive{Color: sp(hexError)},
			GenericEmph:         ansi.StylePrimitive{Italic: bp(true)},
			GenericInserted:     ansi.StylePrimitive{Color: sp("#A3E635")},
			GenericStrong:       ansi.StylePrimitive{Bold: bp(true)},
			GenericSubheading:   ansi.StylePrimitive{Color: sp(hexSubtitle)},
			Background:          ansi.StylePrimitive{BackgroundColor: sp(hexDarkBg)},
		},
		Theme: "monokai",
	},
	Table:                 ansi.StyleTable{CenterSeparator: sp("┼"), ColumnSeparator: sp("│"), RowSeparator: sp("─")},
	DefinitionDescription: ansi.StylePrimitive{BlockPrefix: "\n🠶 "},
}

func AmberStyle() ansi.StyleConfig {
	return amberStyle
}
