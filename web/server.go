// Package web provides the HTTP server that serves minimal static HTML pages
// for each piece of site content. Pages are intentionally plain — no styles,
// no scripts — with a prompt directing visitors to the SSH experience.
package web

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/yuin/goldmark"
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
			markdownToHTML(item.Body),
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
			markdownToHTML(item.Body),
		)
		writePage(w, item.Title+" — "+site.Site.Title, body)
	}
}

// markdownToHTML converts markdown to an HTML string.
// On error it falls back to returning the raw markdown in a <pre> block.
func markdownToHTML(md string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "<pre>" + html.EscapeString(md) + "</pre>"
	}
	return buf.String()
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
