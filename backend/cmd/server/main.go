package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"todo-app/backend/internal/initializer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, cfg, cleanup, err := initializer.BuildServer(ctx)
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	errCh := make(chan error, 1)

	go func() {
		srv.Addr = fmt.Sprintf(":%d", cfg.ServerPort)
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdown(srv, cleanup)
	case err := <-errCh:
		if err != nil {
			cleanup()
			log.Fatalf("server error: %v", err)
		}
	}
}

func shutdown(srv *http.Server, cleanup func()) {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	cleanup()
	log.Println("Server stopped")
}
