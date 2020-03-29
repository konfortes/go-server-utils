package serverutils

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GetEnvOr gets an environment variable or returns ifNotFound value
func GetEnvOr(env, ifNotFound string) string {
	foundEnv, found := os.LookupEnv(env)

	if found {
		return foundEnv
	}

	return ifNotFound
}

func gracefulShutdown(srv *http.Server, shutdownHooks []func()) {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, hook := range shutdownHooks {
		hook()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
