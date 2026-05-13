package components

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

var imageRe = regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)\s*$`)
var attrRe = regexp.MustCompile(`^\{\s*width\s*=\s*"?(\d+)"?\s*\}\s*$`)

// MarkdownContent is a value-type Renderer that renders markdown body text
// with inline image support via half-block art.
type MarkdownContent struct {
	Body       string
	FS         fs.FS
	ContentDir string
}

// imagesReadyMsg is delivered when background image rendering completes.
type imagesReadyMsg struct{ content string }

// Render implements Renderer. Uses AppContext width for word wrapping.
func (mc MarkdownContent) Render(m AppContext) string {
	out, err := mc.renderWidth(m.Width())
	if err != nil {
		return mc.Body
	}
	return out
}

// ViewportFast creates a viewport immediately using text-only rendering.
// Images appear as placeholders until RenderImagesCmd completes.
func (mc MarkdownContent) ViewportFast(width, height int) viewport.Model {
	rendered, err := mc.renderWidthFast(width)
	if err != nil {
		rendered = mc.Body
	}
	vp := viewport.New(width, height)
	vp.SetContent(rendered)
	return vp
}

// RenderImagesCmd returns a Cmd that decodes and renders images in the
// background, delivering imagesReadyMsg when done. Returns nil if there
// are no images to render.
func (mc MarkdownContent) RenderImagesCmd(width int) tea.Cmd {
	if mc.FS == nil {
		return nil
	}
	hasImages := false
	for _, s := range splitAtImages(mc.Body) {
		if s.isImage {
			hasImages = true
			break
		}
	}
	if !hasImages {
		return nil
	}
	return func() tea.Msg {
		rendered, err := mc.renderWidth(width)
		if err != nil {
			return nil
		}
		return imagesReadyMsg{content: rendered}
	}
}

func (mc MarkdownContent) renderWidth(width int) (string, error) {
	w := width - 6
	if w < 20 {
		w = 20
	}

	if mc.FS == nil {
		return renderGlamour(mc.Body, w)
	}

	segments := splitAtImages(mc.Body)
	var parts []string

	for _, seg := range segments {
		if seg.isImage {
			rendered, err := renderImage(mc.FS, mc.ContentDir, seg.imgPath, w, seg.width)
			if err != nil {
				placeholder := "🖼 " + seg.alt
				if seg.alt == "" {
					placeholder = "🖼 " + seg.imgPath
				}
				out, _ := renderGlamour(placeholder, w)
				parts = append(parts, out)
			} else {
				parts = append(parts, rendered)
			}
		} else if strings.TrimSpace(seg.text) != "" {
			out, err := renderGlamour(seg.text, w)
			if err != nil {
				parts = append(parts, seg.text)
			} else {
				parts = append(parts, out)
			}
		}
	}

	return strings.Join(parts, ""), nil
}

// renderWidthFast renders text segments only; images become 🖼 placeholders.
func (mc MarkdownContent) renderWidthFast(width int) (string, error) {
	w := width - 6
	if w < 20 {
		w = 20
	}

	if mc.FS == nil {
		return renderGlamour(mc.Body, w)
	}

	segments := splitAtImages(mc.Body)
	var parts []string

	for _, seg := range segments {
		if seg.isImage {
			placeholder := "🖼 " + seg.alt
			if seg.alt == "" {
				placeholder = "🖼 " + seg.imgPath
			}
			out, _ := renderGlamour(placeholder, w)
			parts = append(parts, out)
		} else if strings.TrimSpace(seg.text) != "" {
			out, err := renderGlamour(seg.text, w)
			if err != nil {
				parts = append(parts, seg.text)
			} else {
				parts = append(parts, out)
			}
		}
	}

	return strings.Join(parts, ""), nil
}

func renderGlamour(content string, wordWrap int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(AmberStyle()),
		glamour.WithWordWrap(wordWrap),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}

type segment struct {
	isImage bool
	text    string
	imgPath string
	alt     string
	width   int // requested display width in columns; 0 = natural size
}

func splitAtImages(content string) []segment {
	lines := strings.Split(content, "\n")
	var segments []segment
	var textBuf []string

	flush := func() {
		if len(textBuf) > 0 {
			segments = append(segments, segment{text: strings.Join(textBuf, "\n")})
			textBuf = nil
		}
	}

	for i := 0; i < len(lines); i++ {
		m := imageRe.FindStringSubmatch(lines[i])
		if m == nil {
			textBuf = append(textBuf, lines[i])
			continue
		}
		flush()
		var w int
		// Look ahead for an attribute line like {width="40"}.
		if i+1 < len(lines) {
			if am := attrRe.FindStringSubmatch(lines[i+1]); am != nil {
				fmt.Sscanf(am[1], "%d", &w)
				i++ // consume the attribute line
			}
		}
		segments = append(segments, segment{
			isImage: true,
			alt:     m[1],
			imgPath: m[2],
			width:   w,
		})
	}
	flush()

	return segments
}
