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

// Theme holds all per-session lipgloss styles and colors. Every style is bound
// to a session-specific *lipgloss.Renderer so adaptive colors resolve against
// the connecting client's terminal background — never a process-global one.
type Theme struct {
	r    *lipgloss.Renderer
	dark bool

	// Raw colors, for building inline styles via NewStyle.
	Primary  lipgloss.TerminalColor
	Accent   lipgloss.TerminalColor
	Muted    lipgloss.TerminalColor
	Normal   lipgloss.TerminalColor
	Clock    lipgloss.TerminalColor
	Subtitle lipgloss.TerminalColor
	Error    lipgloss.TerminalColor

	// Prebuilt styles.
	TitleStyle        lipgloss.Style
	SubtitleStyle     lipgloss.Style
	HeaderStyle       lipgloss.Style
	SelectedItemStyle lipgloss.Style
	NormalItemStyle   lipgloss.Style
	MutedStyle        lipgloss.Style
	ClockStyle        lipgloss.Style
	HelpStyle         lipgloss.Style
	ErrorStyle        lipgloss.Style
	AsciiGlowStyle    lipgloss.Style
	WidgetBorderStyle lipgloss.Style
	WidgetTitleStyle  lipgloss.Style
	WidgetAccentStyle lipgloss.Style
}

// NewTheme builds a Theme bound to r. If r is nil the lipgloss default renderer
// is used (useful for tests). The dark/light decision is taken from the
// renderer's detected background at construction time.
func NewTheme(r *lipgloss.Renderer) *Theme {
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}
	t := &Theme{
		r:        r,
		dark:     r.HasDarkBackground(),
		Primary:  colorPrimary,
		Accent:   colorAccent,
		Muted:    colorMuted,
		Normal:   colorNormal,
		Clock:    colorClock,
		Subtitle: colorSubtitle,
		Error:    colorError,
	}
	t.TitleStyle = r.NewStyle().Foreground(colorPrimary).Bold(true)
	t.SubtitleStyle = r.NewStyle().Foreground(colorSubtitle)
	t.HeaderStyle = r.NewStyle().Foreground(colorPrimary).Bold(true)
	t.SelectedItemStyle = r.NewStyle().Foreground(colorPrimary).Bold(true)
	t.NormalItemStyle = r.NewStyle().Foreground(colorNormal)
	t.MutedStyle = r.NewStyle().Foreground(colorMuted)
	t.ClockStyle = r.NewStyle().Foreground(colorClock)
	t.HelpStyle = r.NewStyle().Foreground(colorMuted)
	t.ErrorStyle = r.NewStyle().Foreground(colorError).Bold(true)
	t.AsciiGlowStyle = r.NewStyle().Foreground(colorAccent).Bold(true)
	t.WidgetBorderStyle = r.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorMuted)
	t.WidgetTitleStyle = r.NewStyle().Foreground(colorMuted).Bold(false)
	t.WidgetAccentStyle = r.NewStyle().Foreground(colorPrimary).Bold(true)
	return t
}

// NewStyle returns a fresh style bound to the session renderer.
func (t *Theme) NewStyle() lipgloss.Style { return t.r.NewStyle() }

// Dark reports whether the session's terminal has a dark background.
func (t *Theme) Dark() bool { return t.dark }

// Markdown returns the glamour style config matching this session's background.
func (t *Theme) Markdown() ansi.StyleConfig { return amberStyleFor(t.dark) }

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

// amberStyleFor returns the markdown style config for the given background.
func amberStyleFor(dark bool) ansi.StyleConfig {
	if dark {
		return buildAmberStyle(darkMarkdownPalette)
	}
	return buildAmberStyle(lightMarkdownPalette)
}
