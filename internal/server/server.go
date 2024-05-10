package server

import (
	"net"
	"ssh-portfolio/internal/sshSession"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"

	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	_ "github.com/joho/godotenv/autoload"
)

func CreateServer(port, host string) (*ssh.Server, error) {
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true // Accept all public keys
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(sshSession.CreateHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)

	return server, err
}
