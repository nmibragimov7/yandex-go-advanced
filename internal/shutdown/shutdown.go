package shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown - function that handles graceful shutdown of the server
func Shutdown(ctx context.Context, server *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		log.Println("context canceled, shutting down...")
	case sig := <-stop:
		log.Printf("received signal: %s, shutting down...", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutdown: %v", err)
	} else {
		log.Println("server shut down gracefully")
	}
}
