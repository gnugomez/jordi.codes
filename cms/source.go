package cms

// ContentSource is the contract for loading content items from a specific source.
type ContentSource interface {
	Load(ct ContentType) ([]ContentItem, error)
	LoadBySlug(ct ContentType, slug string) (ContentItem, error)
}

// LocalSource loads content from markdown files in the filesystem.
type LocalSource struct {
	site *Site
}

func NewLocalSource(site *Site) *LocalSource {
	return &LocalSource{site: site}
}

func (s *LocalSource) Load(ct ContentType) ([]ContentItem, error) {
	return s.site.loadLocalContentItems(ct)
}

func (s *LocalSource) LoadBySlug(ct ContentType, slug string) (ContentItem, error) {
	return s.site.loadLocalContentItemBySlug(ct, slug)
}

// GitHubPinnedSource loads content from GitHub pinned repositories.
type GitHubPinnedSource struct {
	site *Site
}

func NewGitHubPinnedSource(site *Site) *GitHubPinnedSource {
	return &GitHubPinnedSource{site: site}
}

func (s *GitHubPinnedSource) Load(ct ContentType) ([]ContentItem, error) {
	return s.site.loadGitHubPinnedItems(ct)
}

func (s *GitHubPinnedSource) LoadBySlug(ct ContentType, slug string) (ContentItem, error) {
	return s.site.loadGitHubPinnedItemBySlug(ct, slug)
}

// sourceRegistry holds instantiated content sources for the site.
type sourceRegistry struct {
	local        *LocalSource
	githubPinned *GitHubPinnedSource
}

func newSourceRegistry(site *Site) *sourceRegistry {
	return &sourceRegistry{
		local:        NewLocalSource(site),
		githubPinned: NewGitHubPinnedSource(site),
	}
}

func (r *sourceRegistry) Get(sourceType string) ContentSource {
	switch contentTypeSource(sourceType) {
	case "github_pinned":
		return r.githubPinned
	default:
		return r.local
	}
}
