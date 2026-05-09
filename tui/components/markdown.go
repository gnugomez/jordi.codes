package components

import "github.com/charmbracelet/glamour"

func RenderMarkdown(content string, width int) (string, error) {
	w := width - 6
	if w < 20 {
		w = 20
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(AmberStyle()),
		glamour.WithWordWrap(w),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}
