package cms

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ContentTypeByName looks up a ContentType by name.
func (c *Site) ContentTypeByName(name string) *ContentType {
	for i := range c.ContentTypes {
		if c.ContentTypes[i].Name == name {
			return &c.ContentTypes[i]
		}
	}
	return nil
}

// LoadMenuContentItems loads the content list for a menu entry of type content_type.
func (c *Site) LoadMenuContentItems(entry MenuEntry) ([]ContentItem, error) {
	ct := c.ContentTypeByName(entry.ContentType)
	if ct == nil {
		return nil, fmt.Errorf("content type %q not found", entry.ContentType)
	}
	return c.LoadContentItems(*ct)
}

// LoadStaticPage reads a single markdown file and returns its raw content.
func (c *Site) LoadStaticPage(filePath string) (string, error) {
	raw, err := fs.ReadFile(c.fsys, filePath)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// FindContentType looks up a ContentType by name in the given config.
func FindContentType(site *Site, name string) *ContentType {
	return site.ContentTypeByName(name)
}

// LoadContentItems reads items from the appropriate content source based on the ContentType.
func (c *Site) LoadContentItems(ct ContentType) ([]ContentItem, error) {
	source := c.sourceRegistry.Get(ct.Source)
	return source.Load(ct)
}

func (c *Site) loadLocalContentItems(ct ContentType) ([]ContentItem, error) {
	entries, err := fs.ReadDir(c.fsys, ct.Folder)
	if err != nil {
		return nil, err
	}

	var items []ContentItem
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if entry.Name() == "_index.md" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ".md")
		raw, err := fs.ReadFile(c.fsys, path.Join(ct.Folder, entry.Name()))
		if err != nil {
			continue
		}

		meta, body := parseFrontMatter(string(raw))
		title := meta["title"]
		if title == "" {
			title = firstHeading(body, slug)
		}

		items = append(items, ContentItem{
			Title:      title,
			Slug:       slug,
			Path:       "/" + ct.Name + "/" + slug,
			Body:       body,
			Excerpt:    extractExcerpt(body, 200),
			ContentDir: path.Join(ct.Folder, slug),
		})
	}
	return items, nil
}

func contentTypeSource(source string) string {
	src := strings.TrimSpace(strings.ToLower(source))
	if src == "" {
		return "local"
	}
	return src
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
			strings.HasPrefix(line, "![") ||
			strings.HasPrefix(line, "{") ||
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
func (c *Site) LoadStatic(filePath string) (ContentItem, error) {
	raw, err := fs.ReadFile(c.fsys, filePath)
	if err != nil {
		return ContentItem{}, err
	}
	meta, body := parseFrontMatter(string(raw))
	slug := strings.TrimSuffix(path.Base(filePath), ".md")
	title := meta["title"]
	if title == "" {
		title = firstHeading(body, slug)
	}
	return ContentItem{
		Title:      title,
		Slug:       slug,
		Path:       "/" + slug,
		Body:       body,
		Excerpt:    extractExcerpt(body, 200),
		ContentDir: path.Dir(filePath),
	}, nil
}

// LoadContentItemBySlug loads a single item from a ContentType by its slug.
func (c *Site) LoadContentItemBySlug(entry MenuEntry, slug string) (ContentItem, error) {
	ct := c.ContentTypeByName(entry.ContentType)
	if ct == nil {
		return ContentItem{}, fmt.Errorf("content type %q not found", entry.ContentType)
	}
	source := c.sourceRegistry.Get(ct.Source)
	return source.LoadBySlug(*ct, strings.TrimSpace(slug))
}

func (c *Site) loadLocalContentItemBySlug(ct ContentType, slug string) (ContentItem, error) {
	raw, err := fs.ReadFile(c.fsys, path.Join(ct.Folder, slug+".md"))
	if err != nil {
		return ContentItem{}, err
	}
	meta, body := parseFrontMatter(string(raw))
	title := meta["title"]
	if title == "" {
		title = firstHeading(body, slug)
	}
	return ContentItem{
		Title:      title,
		Slug:       slug,
		Path:       "/" + ct.Name + "/" + slug,
		Body:       body,
		Excerpt:    extractExcerpt(body, 200),
		ContentDir: path.Join(ct.Folder, slug),
	}, nil
}
