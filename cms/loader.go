package cms

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// LoadStaticPage reads a single markdown file and returns its raw content.
func LoadStaticPage(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// FindContentType looks up a ContentType by name in the given config.
func FindContentType(cfg *Config, name string) *ContentType {
	for i := range cfg.ContentTypes {
		if cfg.ContentTypes[i].Name == name {
			return &cfg.ContentTypes[i]
		}
	}
	return nil
}

// LoadContentItems reads all markdown files in a ContentType's folder and
// returns them as ContentItem slices. Files are returned in directory order.
func LoadContentItems(ct ContentType) ([]ContentItem, error) {
	entries, err := os.ReadDir(ct.Folder)
	if err != nil {
		return nil, err
	}

	var items []ContentItem
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ".md")
		raw, err := os.ReadFile(filepath.Join(ct.Folder, entry.Name()))
		if err != nil {
			continue
		}

		meta, body := parseFrontMatter(string(raw))
		title := meta["title"]
		if title == "" {
			title = firstHeading(body, slug)
		}

		items = append(items, ContentItem{
			Title:   title,
			Slug:    slug,
			Body:    body,
			Excerpt: extractExcerpt(body, 200),
		})
	}
	return items, nil
}

// parseFrontMatter extracts YAML-like key: value pairs from a leading ---
// block and returns the remaining body content.
func parseFrontMatter(content string) (meta map[string]string, body string) {
	meta = make(map[string]string)
	if !strings.HasPrefix(content, "---\n") {
		return meta, content
	}
	rest := content[4:]
	endIdx := strings.Index(rest, "\n---\n")
	if endIdx == -1 {
		return meta, content
	}
	for _, line := range strings.Split(rest[:endIdx], "\n") {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			meta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return meta, rest[endIdx+5:]
}

// firstHeading returns the text of the first # Heading found in the markdown,
// falling back to a humanised version of the slug.
func firstHeading(body, slug string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return strings.ReplaceAll(cases.Title(language.English).String(slug), "-", " ")
}

// extractExcerpt returns the first meaningful prose line from body, capped at
// maxLen runes. Headings, blank lines, code fences, and table rows are skipped.
func extractExcerpt(body string, maxLen int) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" ||
			strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "```") ||
			strings.HasPrefix(line, "|") ||
			strings.Trim(line, "-=_*~ ") == "" {
			continue
		}
		runes := []rune(line)
		if len(runes) > maxLen {
			return string(runes[:maxLen-1]) + "…"
		}
		return line
	}
	return ""
}

// LoadStatic reads a single markdown file and returns it as a ContentItem.
func LoadStatic(path string) (ContentItem, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ContentItem{}, err
	}
	meta, body := parseFrontMatter(string(raw))
	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	title := meta["title"]
	if title == "" {
		title = firstHeading(body, slug)
	}
	return ContentItem{
		Title:   title,
		Slug:    slug,
		Body:    body,
		Excerpt: extractExcerpt(body, 200),
	}, nil
}
