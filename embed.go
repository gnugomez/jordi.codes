package main

import "embed"

// SiteFS holds the embedded content and config directories.
// All markdown content and site configuration is baked into the binary
// at compile time so no filesystem access is needed at runtime.
//
//go:embed content config
var SiteFS embed.FS

// PublicFS holds the Hugo-generated static site output.
// Run `hugo` before `go build` to populate the public/ directory.
//
//go:embed public
var PublicFS embed.FS
