package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"ssh-portfolio/internal/server"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	appEnv := os.Getenv("APP_ENV")

	server, err := server.CreateServer(port, host)
	if err != nil {
		log.Error("Could not create server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	if appEnv == "local" {
		formattedMsg := fmt.Sprintf("To SSH into the server, run `ssh -p %s %s`", port, host)
		log.Info(formattedMsg)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
