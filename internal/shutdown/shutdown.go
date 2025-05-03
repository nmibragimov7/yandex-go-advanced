package shutdown

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown - function that handles graceful shutdown of the server
func Shutdown(server *http.Server, timeout time.Duration) {
	// Создаем контекст, который будет отменён при получении сигнала завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Ожидаем отмены контекста (например, по сигналу)
	<-ctx.Done()
	log.Println("shutdown signal received, shutting down...")

	// Контекст с таймаутом на завершение сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutdown: %v", err)
	} else {
		log.Println("server shut down gracefully")
	}
}
