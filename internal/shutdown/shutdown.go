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

func GracefulShutdown(server *http.Server, serverErr <-chan error, cancelFuncs ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-serverErr:
		if err != nil {
			log.Printf("HTTP server stopped unexpectedly: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	for _, cancelFunc := range cancelFuncs {
		cancelFunc()
	}

	log.Println("Server exited gracefully")
}
