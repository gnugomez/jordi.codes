package components

import (
	"os"
	"strings"

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

	lightHexPrimary  = "#F97316"
	lightHexAccent   = "#EA580C"
	lightHexMuted    = "#44403C"
	lightHexNormal   = "#292524"
	lightHexClock    = "#9A3412"
	lightHexSubtitle = "#292524"
	lightHexError    = "#B91C1C"
	lightHexBright   = "#7C2D12"
	lightHexAmber    = "#9A3412"
	lightHexCodeBg   = "#F5F5F4"
)

var (
	colorPrimary  = lipgloss.AdaptiveColor{Light: lightHexPrimary, Dark: hexPrimary}
	colorAccent   = lipgloss.AdaptiveColor{Light: lightHexAccent, Dark: hexAccent}
	colorMuted    = lipgloss.AdaptiveColor{Light: lightHexMuted, Dark: hexMuted}
	colorNormal   = lipgloss.AdaptiveColor{Light: lightHexNormal, Dark: hexNormal}
	colorClock    = lipgloss.AdaptiveColor{Light: lightHexClock, Dark: hexClock}
	colorSubtitle = lipgloss.AdaptiveColor{Light: lightHexSubtitle, Dark: hexSubtitle}
	colorError    = lipgloss.AdaptiveColor{Light: lightHexError, Dark: hexError}
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

type markdownPalette struct {
	primary      string
	accent       string
	muted        string
	normal       string
	clock        string
	subtitle     string
	error        string
	bright       string
	amber        string
	bg           string
	codeBg       string
	stringEscape string
	inserted     string
}

var darkMarkdownPalette = markdownPalette{
	primary:      hexPrimary,
	accent:       hexAccent,
	muted:        hexMuted,
	normal:       hexNormal,
	clock:        hexClock,
	subtitle:     hexSubtitle,
	error:        hexError,
	bright:       hexBright,
	amber:        hexAmber,
	bg:           hexDarkBg,
	codeBg:       hexCodeBg,
	stringEscape: "#F59E0B",
	inserted:     "#A3E635",
}

var lightMarkdownPalette = markdownPalette{
	primary:      lightHexAccent,
	accent:       "#C2410C",
	muted:        lightHexMuted,
	normal:       "#1C1917",
	clock:        lightHexClock,
	subtitle:     lightHexSubtitle,
	error:        lightHexError,
	bright:       lightHexBright,
	amber:        lightHexAmber,
	bg:           lightHexNormal,
	codeBg:       lightHexCodeBg,
	stringEscape: "#B45309",
	inserted:     "#3F6212",
}

func buildAmberStyle(p markdownPalette) ansi.StyleConfig {
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: up(2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:  sp(p.clock),
				Italic: bp(true),
			},
			Indent:      up(1),
			IndentToken: sp("│ "),
		},
		List: ansi.StyleList{LevelIndent: 2},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       sp(p.primary),
				Bold:        bp(true),
			},
		},
		H1:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: " ", Suffix: " ", Color: sp(p.bg), BackgroundColor: sp(p.primary), Bold: bp(true), Upper: bp(true)}},
		H2:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "## ", Color: sp(p.accent), Bold: bp(true)}},
		H3:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "### ", Color: sp(p.clock), Bold: bp(true)}},
		H4:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "#### ", Color: sp(p.clock), Bold: bp(true)}},
		H5:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "##### ", Color: sp(p.clock), Bold: bp(true)}},
		H6:             ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: "###### ", Color: sp(p.clock), Bold: bp(true)}},
		Text:           ansi.StylePrimitive{Color: sp(p.normal)},
		Strikethrough:  ansi.StylePrimitive{CrossedOut: bp(true)},
		Emph:           ansi.StylePrimitive{Italic: bp(true), Color: sp(p.clock)},
		Strong:         ansi.StylePrimitive{Bold: bp(true), Color: sp(p.accent)},
		HorizontalRule: ansi.StylePrimitive{Color: sp(p.muted), Format: "\n--------\n"},
		Item:           ansi.StylePrimitive{BlockPrefix: "• "},
		Enumeration:    ansi.StylePrimitive{BlockPrefix: ". "},
		Task:           ansi.StyleTask{Ticked: "[✓] ", Unticked: "[ ] "},
		Link:           ansi.StylePrimitive{Color: sp(p.clock), Underline: bp(true)},
		LinkText:       ansi.StylePrimitive{Color: sp(p.accent), Bold: bp(true)},
		Image:          ansi.StylePrimitive{Color: sp(p.clock), Underline: bp(true)},
		ImageText:      ansi.StylePrimitive{Color: sp(p.accent), Format: "Image: {{.text}} ->"},
		Code:           ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: " ", Suffix: " ", Color: sp(p.amber), BackgroundColor: sp(p.codeBg)}},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Color: sp(p.normal)}, Margin: up(2)},
			Chroma: &ansi.Chroma{
				Text:                ansi.StylePrimitive{Color: sp(p.normal)},
				Comment:             ansi.StylePrimitive{Color: sp(p.muted)},
				CommentPreproc:      ansi.StylePrimitive{Color: sp(p.clock)},
				Keyword:             ansi.StylePrimitive{Color: sp(p.primary)},
				KeywordReserved:     ansi.StylePrimitive{Color: sp(p.primary)},
				KeywordNamespace:    ansi.StylePrimitive{Color: sp(p.accent)},
				KeywordType:         ansi.StylePrimitive{Color: sp(p.clock)},
				Operator:            ansi.StylePrimitive{Color: sp(p.subtitle)},
				Punctuation:         ansi.StylePrimitive{Color: sp(p.subtitle)},
				Name:                ansi.StylePrimitive{Color: sp(p.normal)},
				NameBuiltin:         ansi.StylePrimitive{Color: sp(p.amber)},
				NameTag:             ansi.StylePrimitive{Color: sp(p.primary)},
				NameAttribute:       ansi.StylePrimitive{Color: sp(p.clock)},
				NameClass:           ansi.StylePrimitive{Color: sp(p.amber), Underline: bp(true), Bold: bp(true)},
				NameConstant:        ansi.StylePrimitive{Color: sp(p.amber)},
				NameDecorator:       ansi.StylePrimitive{Color: sp(p.accent)},
				NameFunction:        ansi.StylePrimitive{Color: sp(p.amber)},
				LiteralNumber:       ansi.StylePrimitive{Color: sp(p.clock)},
				LiteralString:       ansi.StylePrimitive{Color: sp(p.bright)},
				LiteralStringEscape: ansi.StylePrimitive{Color: sp(p.stringEscape)},
				GenericDeleted:      ansi.StylePrimitive{Color: sp(p.error)},
				GenericEmph:         ansi.StylePrimitive{Italic: bp(true)},
				GenericInserted:     ansi.StylePrimitive{Color: sp(p.inserted)},
				GenericStrong:       ansi.StylePrimitive{Bold: bp(true)},
				GenericSubheading:   ansi.StylePrimitive{Color: sp(p.subtitle)},
				Background:          ansi.StylePrimitive{BackgroundColor: sp(p.codeBg)},
			},
			Theme: "monokai",
		},
		Table:                 ansi.StyleTable{CenterSeparator: sp("┼"), ColumnSeparator: sp("│"), RowSeparator: sp("─")},
		DefinitionDescription: ansi.StylePrimitive{BlockPrefix: "\n-> "},
	}
}

func init() {
	if forcedDark, ok := forcedDarkMode(); ok {
		lipgloss.SetHasDarkBackground(forcedDark)
	}
}

func forcedDarkMode() (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("JORDI_THEME"))) {
	case "dark":
		return true, true
	case "light":
		return false, true
	default:
		return false, false
	}
}

func AmberStyle() ansi.StyleConfig {
	if forcedDark, ok := forcedDarkMode(); ok {
		if forcedDark {
			return buildAmberStyle(darkMarkdownPalette)
		}
		return buildAmberStyle(lightMarkdownPalette)
	}

	if lipgloss.HasDarkBackground() {
		return buildAmberStyle(darkMarkdownPalette)
	}

	return buildAmberStyle(lightMarkdownPalette)
}
