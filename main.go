package main

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	cssh "github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"

	"jordi.codes/cms"
	"jordi.codes/router"
	"jordi.codes/tui"
)

const (
	host     = "0.0.0.0"
	sshPort  = "22"
	httpPort = "80"
	keyPath  = ".ssh/id_ed25519"
)

func main() {
	// stdout is not a terminal when running in Docker, so termenv defaults to
	// no-color. Force TrueColor so lipgloss styles are always rendered.
	lipgloss.SetColorProfile(termenv.TrueColor)

	site, err := cms.LoadSite(SiteFS, "config/content.yaml")
	if err != nil {
		log.Fatal("Failed to load site config", "error", err)
	}
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, sshPort)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			// Render the bubbletea TUI for each SSH session.
			bm.Middleware(func(s cssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, active := s.Pty()
				if !active {
					return nil, nil
				}
				addr := s.RemoteAddr().String()
				if h, _, err := net.SplitHostPort(addr); err == nil {
					addr = h
				}

				// Resolve the SSH username to a route path.
				// e.g. "about" → "/about", "projects.jordi-codes" → "/projects/jordi-codes"
				// If the username does not match a known route the TUI starts at the menu.
				initialPath := router.UsernameToPath(s.User())

				m := tui.NewModel(site, pty.Window.Width, pty.Window.Height, addr,
					tui.WithInitialPath(initialPath),
				)
				return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
			}),
			// Reject connections that do not have an active PTY.
			activeterm.Middleware(),
			// Structured request logging via charmbracelet/log.
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("Could not create SSH server", "error", err)
	}

	publicFS, err := fs.Sub(PublicFS, "public")
	if err != nil {
		log.Fatal("Failed to open public directory", "error", err)
	}
	httpSrv := &http.Server{
		Addr:    net.JoinHostPort(host, httpPort),
		Handler: http.FileServer(http.FS(publicFS)),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("SSH server started", "host", host, "port", sshPort)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, cssh.ErrServerClosed) {
			log.Fatal("SSH server error", "error", err)
			done <- nil
		}
	}()

	log.Info("HTTP server started", "host", host, "port", httpPort)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", "error", err)
		}
	}()

	<-done

	log.Info("Shutting down servers…")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, cssh.ErrServerClosed) {
		log.Fatal("SSH graceful shutdown failed", "error", err)
	}
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error("HTTP graceful shutdown failed", "error", err)
	}
}
