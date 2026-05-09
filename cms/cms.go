package cms

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level site configuration loaded from content.yaml.
type Config struct {
	Site         SiteConfig    `yaml:"site"`
	Menu         []MenuEntry   `yaml:"menu"`
	ContentTypes []ContentType `yaml:"content_types"`
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
	Title   string
	Slug    string
	Body    string
	Excerpt string // first meaningful paragraph, plain text
}

// LoadConfig reads and parses the YAML site configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
