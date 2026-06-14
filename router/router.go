// Package router resolves URL paths to site content.
// Both the SSH TUI and the HTTP server use this package to navigate content.
package router

import (
	"path"
	"strings"

	"jordi.codes/cms"
)

// Kind describes the type of content a route points to.
type Kind string

const (
	KindList   Kind = "list"   // a listing of content items (e.g. /projects)
	KindStatic Kind = "static" // a single static markdown page (e.g. /about)
	KindDetail Kind = "detail" // one item within a content type (e.g. /projects/jordi-codes)
)

// Route is a resolved site location with all the information needed to
// render it, whether over SSH or HTTP.
type Route struct {
	Path  string        // canonical URL path, e.g. "/about", "/projects/jordi-codes"
	Title string        // display title for the page
	Kind  Kind          // what kind of content this route points to
	Entry cms.MenuEntry // the underlying menu entry driving this route
	Slug  string        // populated only for KindDetail: the content item slug
}

// Router maps URL paths to site content using the site's menu configuration.
type Router struct {
	site *cms.Site
}

// New creates a Router for the given site.
func New(site *cms.Site) *Router {
	return &Router{site: site}
}

// SlugFromFilePath derives a URL slug from a file path.
// "content/about.md" → "about"
func SlugFromFilePath(filePath string) string {
	return strings.TrimSuffix(path.Base(filePath), ".md")
}

// UsernameToPath converts an SSH username to a URL path using dot-notation.
// "about" → "/about"
// "projects.jordi-codes" → "/projects/jordi-codes"
func UsernameToPath(username string) string {
	return "/" + strings.ReplaceAll(username, ".", "/")
}

// Resolve maps a URL path to a Route based on the site menu config.
// Returns false if no route matches.
func (r *Router) Resolve(urlPath string) (Route, bool) {
	urlPath = strings.TrimPrefix(urlPath, "/")
	if urlPath == "" {
		return Route{}, false
	}

	parts := strings.SplitN(urlPath, "/", 2)
	first := parts[0]

	for _, entry := range r.site.Menu {
		entry := entry // capture loop variable
		switch entry.Type {
		case "static":
			slug := SlugFromFilePath(entry.Path)
			if first == slug && len(parts) == 1 {
				return Route{
					Path:  "/" + slug,
					Title: entry.Label,
					Kind:  KindStatic,
					Entry: entry,
				}, true
			}

		case "content_type":
			ct := r.site.ContentTypeByName(entry.ContentType)
			if ct == nil || first != ct.Name {
				continue
			}
			if len(parts) == 2 && parts[1] != "" {
				return Route{
					Path:  "/" + ct.Name + "/" + parts[1],
					Title: parts[1],
					Kind:  KindDetail,
					Entry: entry,
					Slug:  parts[1],
				}, true
			}
			return Route{
				Path:  "/" + ct.Name,
				Title: entry.Label,
				Kind:  KindList,
				Entry: entry,
			}, true
		}
	}

	return Route{}, false
}
