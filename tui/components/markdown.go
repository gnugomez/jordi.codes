package components

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"jordi.codes/cms"
)

var imageRe = regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)\s*$`)
var attrRe = regexp.MustCompile(`^\{\s*width\s*=\s*"?(\d+)(%)?"?\s*\}\s*$`)
var shortcodeRe = regexp.MustCompile(`^\{\{<\s*([a-zA-Z0-9_-]+)(.*?)\s*>\}\}$`)
var shortcodeArgRe = regexp.MustCompile(`([a-zA-Z0-9_-]+)\s*=\s*"([^"]*)"`)

// MarkdownContent is a value-type Renderer that renders markdown body text
// with inline image support via half-block art.
type MarkdownContent struct {
	Body       string
	FS         fs.FS
	ContentDir string
	Dark       bool
}

// imagesReadyMsg is delivered when background image rendering completes.
type imagesReadyMsg struct{ content string }

// Render implements Renderer. Uses AppContext width for word wrapping.
func (mc MarkdownContent) Render(m AppContext) string {
	out, err := mc.renderWidthWithContext(m.Width(), m)
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
	return mc.renderSegments(width, false, nil)
}

// renderWidthWithContext renders with AppContext available for shortcodes.
func (mc MarkdownContent) renderWidthWithContext(width int, ctx AppContext) (string, error) {
	return mc.renderSegments(width, false, ctx)
}

// renderWidthFast renders text segments only; images become 🖼 placeholders.
func (mc MarkdownContent) renderWidthFast(width int) (string, error) {
	return mc.renderSegments(width, true, nil)
}

func (mc MarkdownContent) renderSegments(width int, placeholdersOnly bool, ctx AppContext) (string, error) {
	body := mc.expandShortcodes(mc.Body, ctx)
	w := width - 6
	if w < 20 {
		w = 20
	}

	if mc.FS == nil {
		return renderGlamour(body, w, mc.Dark)
	}

	segments := splitAtImages(body)
	var parts []string

	for _, seg := range segments {
		if seg.isImage {
			if placeholdersOnly {
				parts = append(parts, renderPlaceholder(w, seg, nil, mc.Dark))
				continue
			}

			effectiveCols := w
			if seg.widthPct > 0 {
				// w = width-6 accounts for glamour text margins, but the 2-char pad
				// in renderHalfBlocks shifts the image right. Use w+2 (= width-4) as
				// the base so 100% gives symmetric 2-char margins on both sides.
				effectiveCols = (w + 2) * seg.widthPct / 100
				if effectiveCols < 1 {
					effectiveCols = 1
				}
			}

			pixelW := seg.width
			if seg.widthPct > 0 {
				pixelW = 0 // percentage controls the column cap; skip pixel scaling
			}

			// When an explicit size is requested, use a large row cap so the
			// height constraint never reduces the width below what was asked for.
			imgMaxRows := 0
			if seg.widthPct > 0 || seg.width > 0 {
				imgMaxRows = 200
			}

			rendered, err := renderImage(mc.FS, mc.ContentDir, seg.imgPath, effectiveCols, pixelW, imgMaxRows)
			if err != nil {
				parts = append(parts, renderPlaceholder(w, seg, err, mc.Dark))
				continue
			}
			parts = append(parts, rendered)
		} else if strings.TrimSpace(seg.text) != "" {
			out, err := renderGlamour(seg.text, w, mc.Dark)
			if err != nil {
				parts = append(parts, seg.text)
			} else {
				parts = append(parts, out)
			}
		}
	}

	return strings.Join(parts, ""), nil
}

func (mc MarkdownContent) expandShortcodes(content string, ctx AppContext) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		m := shortcodeRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(m[1]))
		args := parseShortcodeArgs(m[2])
		rendered, err := mc.renderShortcode(name, args, ctx)
		if err != nil {
			lines[i] = fmt.Sprintf("shortcode %q failed: %v", name, err)
			continue
		}
		lines[i] = rendered
	}
	return strings.Join(lines, "\n")
}

func parseShortcodeArgs(raw string) map[string]string {
	args := map[string]string{}
	for _, m := range shortcodeArgRe.FindAllStringSubmatch(raw, -1) {
		if len(m) == 3 {
			args[strings.ToLower(strings.TrimSpace(m[1]))] = strings.TrimSpace(m[2])
		}
	}
	return args
}

func (mc MarkdownContent) renderShortcode(name string, args map[string]string, ctx AppContext) (string, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case "projects":
		return renderProjectsShortcode(args)
	case "menu":
		return renderMenuShortcode(ctx)
	default:
		return "", fmt.Errorf("unknown shortcode: %s", name)
	}
}

func renderProjectsShortcode(args map[string]string) (string, error) {
	username := strings.TrimSpace(args["user"])
	if username == "" {
		username = strings.TrimSpace(args["github_user"])
	}
	if username == "" {
		username = strings.TrimSpace(os.Getenv("GITHUB_USER"))
	}
	if username == "" {
		return "", fmt.Errorf("projects shortcode requires user=\"...\"")
	}

	limit := 6
	if rawLimit := strings.TrimSpace(args["limit"]); rawLimit != "" {
		if parsed, err := strconv.Atoi(rawLimit); err == nil {
			limit = parsed
		}
	}

	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN is required")
	}

	repos, err := cms.FetchPinnedRepos(username, token, limit)
	if err != nil {
		return "", err
	}
	if len(repos) == 0 {
		return "No pinned repositories found.", nil
	}

	var sb strings.Builder
	sb.WriteString("## Pinned Projects\n\n")
	for _, repo := range repos {
		sb.WriteString("- [")
		sb.WriteString(repo.Name)
		sb.WriteString("](")
		sb.WriteString(repo.URL)
		sb.WriteString(")")
		if strings.TrimSpace(repo.Description) != "" {
			sb.WriteString(" - ")
			sb.WriteString(repo.Description)
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func renderMenuShortcode(ctx AppContext) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("menu shortcode requires AppContext")
	}
	menu := ctx.Menu()
	if len(menu) == 0 {
		return "", fmt.Errorf("menu is empty")
	}

	var sb strings.Builder
	sb.WriteString("## Menu\n\n")
	for _, entry := range menu {
		// Generate SSH route from content type or slug (e.g., "projects" or "about")
		var sshRoute string
		if entry.Type == "content_type" {
			sshRoute = entry.ContentType
		} else {
			// For static pages, extract slug from path (e.g., "content/about.md" → "about")
			sshRoute = strings.TrimSuffix(path.Base(entry.Path), ".md")
		}

		// For nested paths, convert "/" to "." (e.g., "projects/jordi-codes" → "projects.jordi-codes")
		sshRoute = strings.ReplaceAll(sshRoute, "/", ".")

		sb.WriteString("- ")
		sb.WriteString(entry.Label)
		sb.WriteString(" (`ssh ")
		sb.WriteString(sshRoute)
		sb.WriteString("@jordi.codes`)")
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func renderPlaceholder(wordWrap int, seg segment, renderErr error, dark bool) string {
	label := seg.alt
	if label == "" {
		label = seg.imgPath
	}
	placeholder := "🖼 " + label
	if renderErr != nil {
		log.Printf("image render fallback for %q: %v", seg.imgPath, renderErr)
		placeholder += " (image load failed)"
	}
	out, err := renderGlamour(placeholder, wordWrap, dark)
	if err != nil {
		return placeholder + "\n"
	}
	return out
}

func renderGlamour(content string, wordWrap int, dark bool) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(amberStyleFor(dark)),
		glamour.WithWordWrap(wordWrap),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}

type segment struct {
	isImage  bool
	text     string
	imgPath  string
	alt      string
	width    int // requested pixel width; 0 = natural size
	widthPct int // percentage of available columns (1-100); 0 = not set
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
		var pixelW, pctW int
		// Look ahead for an attribute line like {width="400"} or {width="50%"}.
		if i+1 < len(lines) {
			if am := attrRe.FindStringSubmatch(lines[i+1]); am != nil {
				if am[2] == "%" {
					fmt.Sscanf(am[1], "%d", &pctW)
				} else {
					fmt.Sscanf(am[1], "%d", &pixelW)
				}
				i++ // consume the attribute line
			}
		}
		segments = append(segments, segment{
			isImage:  true,
			alt:      m[1],
			imgPath:  m[2],
			width:    pixelW,
			widthPct: pctW,
		})
	}
	flush()

	return segments
}
