package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	host := envOrDefault("HOST", "127.0.0.1")
	port := envOrDefault("PORT", "3001")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Hello from backend-go"}`))
	})

	server := &http.Server{Addr: addr}
	errCh := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	log.Printf("[ ready ] http://%s", addr)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %s, shutting down", sig)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server close error: %v", err)
		}
		return
	case err := <-errCh:
		log.Fatalf("server failed: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
