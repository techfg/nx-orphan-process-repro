package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(&cobra.Command{
		Use:          "serve",
		Short:        "Start Webserver",
		RunE:         serve,
		SilenceUsage: true,
	})

}

var gracefulShutdownSeconds = 5

func serve(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	slog.InfoContext(ctx, "Starting platform server")
	baseRouter := mux.NewRouter()

	port := "3000"
	host := ""
	serveAddr := host + ":" + port

	server := NewServer(serveAddr, baseRouter)
	var serveErr error

	done := make(chan bool)
	go func() {
		slog.InfoContext(ctx, "Service started on port: "+port)
		serveErr = server.ListenAndServe()
		if serveErr != nil && serveErr.Error() != "http: Server closed" {
			slog.ErrorContext(ctx, "failed to start server: "+serveErr.Error())
			// this will terminate the server without waiting for graceful shutdown
			server.StartupError()
		}
		done <- true
	}()

	// wait for graceful shutdown to complete
	server.WaitShutdown()

	<-done

	return nil
}

func NewServer(serveAddress string, router http.Handler) *ServerWithShutdown {
	return &ServerWithShutdown{
		Server: http.Server{
			Addr:    serveAddress,
			Handler: router,
		},
		startupError: make(chan bool),
	}
}

type ServerWithShutdown struct {
	http.Server
	startupError chan bool
}

// StartupError should ONLY be called on initial startup if the server failed to start due to misconfiguration,
// conflicting processes, etc. this will immediately terminate the server,
// without performing graceful shutdown.
func (s *ServerWithShutdown) StartupError() {
	s.startupError <- true
}

// WaitShutdown will block until a SIGINT/SIGTERM signal is sent to the process,
// at which point it will wait for <gracefulShutdownSeconds> before actually killing the server.
// This allows in-flight requests and processes to be cleanly completed.
// Load balancers should NOT be sending us any net-new API requests after the interrupt signal is received.
func (s *ServerWithShutdown) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	startupError := false

	ctx := context.Background()

	// Wait for an interrupt signal is sent
	select {
	case sig := <-irqSig:
		slog.InfoContext(ctx, fmt.Sprintf("Shutdown initiated (signal: %s)", sig.String()))
	case <-s.startupError:
		startupError = true
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(gracefulShutdownSeconds+1)*time.Second)
	defer cancel()

	// If there was an error starting up the server, this channel will receive a message
	if !startupError {
		slog.InfoContext(ctx, fmt.Sprintf("Waiting %d seconds to allow in-flight processes to finish...", gracefulShutdownSeconds))

		// Create shutdown context with timeout
		t := time.NewTimer(time.Duration(gracefulShutdownSeconds) * time.Second)
		defer t.Stop()
		<-t.C
	}

	// Completely shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error terminating server: "+err.Error())
	}
	if !startupError {
		slog.InfoContext(ctx, "Graceful shutdown is complete.")
	}
}
