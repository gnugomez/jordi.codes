package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	cssh "github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"

	"jordi.codes/cms"
	"jordi.codes/tui"
)

const (
	host    = "0.0.0.0"
	port    = "22"
	keyPath = ".ssh/id_ed25519"
)

func main() {
	site, err := cms.LoadSite(SiteFS, "config/content.yaml")
	if err != nil {
		log.Fatal("Failed to load site config", "error", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			// Render the bubbletea TUI for each SSH session.
			bm.Middleware(func(s cssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, active := s.Pty()
				if !active {
					return nil, nil
				}
				addr := s.RemoteAddr().String()
				if host, _, err := net.SplitHostPort(addr); err == nil {
					addr = host
				}
				m := tui.NewModel(site, pty.Window.Width, pty.Window.Height, addr)
				return m, []tea.ProgramOption{tea.WithAltScreen()}
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("SSH server started", "host", host, "port", port)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, cssh.ErrServerClosed) {
			log.Fatal("Server error", "error", err)
			done <- nil
		}
	}()

	<-done

	log.Info("Shutting down SSH server…")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, cssh.ErrServerClosed) {
		log.Fatal("Graceful shutdown failed", "error", err)
	}
}
