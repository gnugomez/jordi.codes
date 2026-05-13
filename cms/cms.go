package cms

import (
	"io/fs"

	"gopkg.in/yaml.v3"
)

// Site is the top-level site configuration loaded from content.yaml.
type Site struct {
	Site         SiteConfig    `yaml:"site"`
	Menu         []MenuEntry   `yaml:"menu"`
	ContentTypes []ContentType `yaml:"content_types"`

	fsys fs.FS // embedded at load time; unexported so callers never touch it
}

// SiteConfig holds global site metadata.
type SiteConfig struct {
	Title    string `yaml:"title"`
	Subtitle string `yaml:"subtitle"`
	Author   string `yaml:"author"`
	GitHub   string `yaml:"github"`
}

// MenuEntry defines one item in the navigation menu.
// Type is either "content_type" (links to a content listing) or "static"
// (renders a single markdown file).
type MenuEntry struct {
	Label       string `yaml:"label"`
	Type        string `yaml:"type"`
	ContentType string `yaml:"content_type"` // used when Type == "content_type"
	Path        string `yaml:"path"`         // used when Type == "static"
}

// ContentType defines a collection of markdown content (e.g. posts, projects).
type ContentType struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"display_name"`
	Folder      string `yaml:"folder"`
}

// ContentItem represents a single markdown entry within a ContentType.
type ContentItem struct {
	Title      string
	Slug       string
	Path       string // canonical URL path, e.g. "/projects/jordi-codes"
	Body       string
	Excerpt    string // first meaningful paragraph, plain text
	ContentDir string // filesystem directory containing the markdown file
}

// LoadSite reads and parses the YAML site configuration file from the given FS.
// The FS is stored inside Config so that content loaders need no external FS argument.
func LoadSite(fsys fs.FS, path string) (*Site, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}
	var site Site
	if err := yaml.Unmarshal(data, &site); err != nil {
		return nil, err
	}
	site.fsys = fsys
	return &site, nil
}

// FS returns the filesystem backing the site content.
func (c *Site) FS() fs.FS { return c.fsys }
