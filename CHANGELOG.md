# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- SSH-served TUI personal site built with Go and the Charmbracelet stack (bubbletea, lipgloss, glamour, wish)
- Menu view with navigation
- GitHub contribution graph widget (requires `GITHUB_TOKEN`)
- Calendar and visitor widgets
- List view with stacked box card layout for browsing content types
- Detail view with scrollable markdown rendering via glamour (amber theme)
- CMS layer: YAML-driven config, front matter parsing, excerpt extraction
- Component-based architecture: `ViewComponent` and `Renderer` interfaces backed by `AppContext`
- View stack navigation (push/pop) with nav message pattern
- Clock in the footer, ticking every second
