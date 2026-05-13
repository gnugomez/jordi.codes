// Package web provides the HTTP server that serves minimal static HTML pages
// for each piece of site content. Pages are intentionally plain — no styles,
// no scripts — with a prompt directing visitors to the SSH experience.
package web

import (
	"bytes"
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"strings"

	"github.com/yuin/goldmark"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"jordi.codes/cms"
	"jordi.codes/router"
)

// NewServer returns an *http.Server that serves plain HTML pages for every
// route in the site. The SSH hostname shown in prompts is derived from the
// HTTP Host header of each incoming request.
func NewServer(site *cms.Site, addr string) *http.Server {
	r := router.New(site)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		sshHost := req.Host
		if req.URL.Path == "/" {
			serveHome(w, site, sshHost)
			return
		}

		// Serve static assets (images) from the content directory.
		if serveContentFile(w, req, site) {
			return
		}

		route, ok := r.Resolve(req.URL.Path)
		if !ok {
			http.NotFound(w, req)
			return
		}
		serveRoute(w, req, route, site, sshHost)
	})

	return &http.Server{Addr: addr, Handler: mux}
}

func serveHome(w http.ResponseWriter, site *cms.Site, sshHost string) {
	var links strings.Builder
	for _, entry := range site.Menu {
		entry := entry
		var href string
		switch entry.Type {
		case "static":
			href = "/" + router.SlugFromFilePath(entry.Path)
		case "content_type":
			ct := cms.FindContentType(site, entry.ContentType)
			if ct == nil {
				continue
			}
			href = "/" + ct.Name
		default:
			continue
		}
		fmt.Fprintf(&links, "<li><a href=%q>%s</a></li>\n", href, html.EscapeString(entry.Label))
	}

	body := fmt.Sprintf(
		"<h1>%s</h1>\n<p>%s</p>\n<ul>\n%s</ul>\n<hr>\n<p>This site is meant to be accessed over SSH:</p>\n<pre><code style=\"border: 1px solid; padding: 0.5em 1em;\">ssh %s</code></pre>\n",
		html.EscapeString(site.Site.Title),
		html.EscapeString(site.Site.Subtitle),
		links.String(),
		html.EscapeString(sshHost),
	)
	writePage(w, site.Site.Title, body)
}

func serveRoute(w http.ResponseWriter, req *http.Request, route router.Route, site *cms.Site, sshHost string) {
	switch route.Kind {
	case router.KindStatic:
		item, err := site.LoadStatic(route.Entry.Path)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		sshCmd := sshLoginCmd(sshHost, router.SlugFromFilePath(route.Entry.Path))
		body := fmt.Sprintf(
			"<p><a href=\"/\">&larr; Back</a></p>\n<h1>%s</h1>\n<p>This site is SSH-first. Connect to read this page:</p>\n<pre><code style=\"border: 1px solid; padding: 0.5em 1em;\">%s</code></pre>\n<hr>\n%s",
			html.EscapeString(item.Title),
			html.EscapeString(sshCmd),
			markdownToHTML(item.Body, item.ContentDir),
		)
		writePage(w, item.Title+" — "+site.Site.Title, body)

	case router.KindList:
		ct := cms.FindContentType(site, route.Entry.ContentType)
		if ct == nil {
			http.NotFound(w, req)
			return
		}
		items, err := site.LoadContentItems(*ct)
		if err != nil {
			http.Error(w, "failed to load content", http.StatusInternalServerError)
			return
		}
		var list strings.Builder
		for _, item := range items {
			fmt.Fprintf(&list, "<li><a href=%q>%s</a></li>\n",
				route.Path+"/"+item.Slug, html.EscapeString(item.Title))
		}
		body := fmt.Sprintf(
			"<p><a href=\"/\">&larr; Back</a></p>\n<h1>%s</h1>\n<ul>\n%s</ul>\n<hr>\n<p>This site is SSH-first. Connect to browse:</p>\n<pre><code style=\"border: 1px solid; padding: 0.5em 1em;\">ssh %s</code></pre>\n",
			html.EscapeString(route.Title),
			list.String(),
			html.EscapeString(sshHost),
		)
		writePage(w, route.Title+" — "+site.Site.Title, body)

	case router.KindDetail:
		item, err := site.LoadContentItemBySlug(route.Entry, route.Slug)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		// Convert "/projects/jordi-codes" → "projects.jordi-codes" for SSH username
		sshUser := strings.ReplaceAll(strings.TrimPrefix(route.Path, "/"), "/", ".")
		sshCmd := sshLoginCmd(sshHost, sshUser)
		parentPath := route.Path[:strings.LastIndex(route.Path, "/")]
		if parentPath == "" {
			parentPath = "/"
		}
		body := fmt.Sprintf(
			"<p><a href=%q>&larr; Back</a></p>\n<h1>%s</h1>\n<p>This site is SSH-first. Connect to read this page:</p>\n<pre><code style=\"border: 1px solid; padding: 0.5em 1em;\">%s</code></pre>\n<hr>\n%s",
			parentPath,
			html.EscapeString(item.Title),
			html.EscapeString(sshCmd),
			markdownToHTML(item.Body, item.ContentDir),
		)
		writePage(w, item.Title+" — "+site.Site.Title, body)
	}
}

// markdownToHTML converts markdown to an HTML string.
// contentDir is the content directory (e.g. "content/recipes") used to resolve
// relative image paths. On error it falls back to the raw markdown in a <pre> block.
func markdownToHTML(md, contentDir string) string {
	// Rewrite relative image paths so they resolve from the web root.
	// e.g. ![](images/foo.jpg) in content/recipes → ![](/recipes/images/foo.jpg)
	md = rewriteImagePaths(md, contentDir)

	gm := goldmark.New(goldmark.WithRendererOptions(goldhtml.WithUnsafe()))
	var buf bytes.Buffer
	if err := gm.Convert([]byte(md), &buf); err != nil {
		return "<pre>" + html.EscapeString(md) + "</pre>"
	}
	// Add max-width to all images.
	result := strings.ReplaceAll(buf.String(), "<img ", "<img style=\"max-width:100%;height:auto\" ")
	return result
}

// rewriteImagePaths rewrites relative image paths in markdown so they resolve
// from the web root. contentDir like "content/recipes" becomes "/recipes".
func rewriteImagePaths(md, contentDir string) string {
	if contentDir == "" {
		return md
	}
	webDir := "/" + strings.TrimPrefix(contentDir, "content/")

	lines := strings.Split(md, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "![") {
			continue
		}
		// Find the src in ![alt](src)
		start := strings.Index(trimmed, "](")
		if start == -1 {
			continue
		}
		end := strings.Index(trimmed[start+2:], ")")
		if end == -1 {
			continue
		}
		src := trimmed[start+2 : start+2+end]
		if strings.HasPrefix(src, "http") || strings.HasPrefix(src, "/") {
			continue
		}
		newSrc := webDir + "/" + src
		lines[i] = strings.Replace(line, "]("+src+")", "]("+newSrc+")", 1)
	}
	return strings.Join(lines, "\n")
}

// serveContentFile attempts to serve a static file (image, etc.) from the
// embedded content directory. Returns true if the file was found and served.
func serveContentFile(w http.ResponseWriter, req *http.Request, site *cms.Site) bool {
	// Map URL path like /recipes/images/foo.jpg → content/recipes/images/foo.jpg
	p := strings.TrimPrefix(req.URL.Path, "/")
	fsPath := "content/" + p

	f, err := site.FS().Open(fsPath)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.(interface{ Stat() (fs.FileInfo, error) }).Stat()
	if err != nil || stat.IsDir() {
		return false
	}

	// Determine content type from extension.
	switch {
	case strings.HasSuffix(p, ".jpg"), strings.HasSuffix(p, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(p, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(p, ".gif"):
		w.Header().Set("Content-Type", "image/gif")
	case strings.HasSuffix(p, ".webp"):
		w.Header().Set("Content-Type", "image/webp")
	case strings.HasSuffix(p, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		return false // only serve known image types
	}

	data, err := fs.ReadFile(site.FS(), fsPath)
	if err != nil {
		return false
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(data)
	return true
}

// sshLoginCmd builds the SSH command string for a given host and username.
func sshLoginCmd(host, username string) string {
	return fmt.Sprintf("ssh %s@%s", username, host)
}

func writePage(w http.ResponseWriter, title, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w,
		"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n<title>%s</title>\n</head>\n<body>\n%s</body>\n</html>",
		html.EscapeString(title),
		body,
	)
}
