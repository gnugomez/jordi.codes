package main

import "embed"

// SiteFS holds the embedded content and config directories.
// All markdown content and site configuration is baked into the binary
// at compile time so no filesystem access is needed at runtime.
//
//go:embed content config
var SiteFS embed.FS
