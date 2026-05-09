# AGENTS.md

Agent instructions for `jordi.codes` — an SSH-served TUI personal website written in Go.

---

## Build & Run

```sh
# Verify the build compiles cleanly (run this after every change)
go build ./...

# Run the server locally (requires the host key and config to be present)
go run .
```

---

## Core Interfaces (`tui/components/contract.go`)

### `AppContext`
Read-only view of application state. Components receive a pointer (`*tui.Context`) that satisfies this interface — they must **never** type-assert it back to a concrete type.

```go
type AppContext interface {
    Width() int
    Height() int
    SetHelpText(string)
    HelpText() string
    Menu() []cms.MenuEntry
    Subtitle() string
    Now() time.Time
    Contribs() map[string]int
    RemoteAddr() string
    ListLayout() ListLayout
}
```

### `Renderer`
A value type (usually a struct) that draws itself given an `AppContext`.

```go
type Renderer interface {
    Render(AppContext) string
}
```

### `ViewComponent`
A stateful pointer type that can both render and handle messages.

```go
type ViewComponent interface {
    Renderer
    Update(m AppContext, msg tea.Msg) tea.Cmd
}
```

---

## Navigation Pattern

  Components **never** mutate the view stack directly. Instead they return a `tea.Cmd` that emits a nav message. `tui.Context.Update()` (in `app.go`) handles the message and calls `pushView` / `popView`.

  | Component wants to…      | Returns                                        | Handled in `app.go` as…      |
  |--------------------------|------------------------------------------------|-------------------------------|
  | Open a content list      | `RequestOpenContentType(entry)` → `OpenContentTypeMsg` | `pushView(NewListView(...))`  |
  | Open a static page       | `RequestOpenStaticPage(entry)` → `OpenStaticPageMsg`   | `pushView(NewStaticDetailView(...))` |
  | Open a detail page       | `RequestOpenDetail(item)` → `OpenDetailMsg`            | `pushView(NewDetailView(...))`|
  | Go back                  | `RequestNavBack()` → `NavBackMsg`              | `popView()`                   |

---

## Component Conventions

1. **Struct-based Renderers.** Use a struct implementing `Render(AppContext) string`. Do not use free functions for rendering.
2. **Components own their state.** `NewListView`, `NewDetailView`, `NewStaticDetailView` load CMS content in their constructors. They store it internally and never require the caller to pass pre-loaded data.
3. **ViewComponents are pointer receivers.** `MenuView`, `ListView`, `DetailView` implement `ViewComponent` with pointer receivers so they can mutate cursor/scroll state in `Update`.
4. **Renderers are value receivers.** `Footer`, `SectionHeader`, `ErrorView`, `ListBody`, `DetailBody`, `MenuBody` are all value types.
5. **No direct tui/ imports from components/.** The `AppContext` interface in `components` breaks the import cycle — `components` knows nothing about `tui`.
6. **`SetHelpText` is called in `Render`, not `Update`.** Each view sets the footer help text during its own render pass.

---

## Styling

  All colors and lipgloss styles live in `tui/components/theme.go`. Use the exported/unexported variables from that file rather than defining styles inline. The amber color palette is used throughout. Markdown is rendered with `glamour` using `AmberStyle()`.

---

## Content Configuration

`config/content.yaml` drives the menu and content types. A `MenuEntry` has:
- `type: content_type` — loads items from a directory via `LoadContentItems`
- `type: static` — loads a single Markdown file via `LoadStaticPage`

Content files live under `content/`. Front matter (YAML) provides title and date; the first `#` heading is used as a fallback title.

---

## Key Dependencies

  | Package | Purpose |
  |---------|---------|
  | `charmbracelet/wish` | SSH server scaffolding |
  | `charmbracelet/bubbletea` | Elm-architecture TUI framework |
  | `charmbracelet/lipgloss` | Terminal styling |
  | `charmbracelet/glamour` | Markdown → ANSI rendering |
  | `charmbracelet/bubbles/viewport` | Scrollable viewport widget |
  | `gopkg.in/yaml.v3` | Config and front matter parsing |

---

## What to Avoid

- Do **not** type-assert `AppContext` back to `*tui.Context` inside `components/`.
- Do **not** add new free-function renderers — add a struct implementing `Renderer` instead.
- Do **not** have components push/pop views themselves — emit a nav message and let `app.go` handle it.
- Do **not** import `jordi.codes/tui` from `jordi.codes/tui/components` or `jordi.codes/tui/layout` (circular).
- Do **not** import `jordi.codes/tui/components` from `jordi.codes/tui/layout` (layout is a leaf package).
