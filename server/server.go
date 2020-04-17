package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konfortes/go-server-utils/logging"
	"github.com/konfortes/go-server-utils/monitoring"
	"github.com/konfortes/go-server-utils/tracing"
)

// Config is server configuration
type Config struct {
	AppName       string
	Port          string
	Env           string
	Handlers      []Handler
	ShutdownHooks []func()
	WithTracing   bool
}

// Handler is a gin HandlerFunc with the http method and pattern to handle
type Handler struct {
	Method  string
	Pattern string
	H       gin.HandlerFunc
}

// Initialize initializes a gin app
func Initialize(config Config) *http.Server {
	router := gin.New()

	var logger gin.HandlerFunc
	if config.Env == "production" {
		logger = logging.JSONLogMiddleware()
	} else {
		logger = gin.Logger()
	}

	router.Use(logger, gin.Recovery(), logging.RequestIDMiddleware())

	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	monitoring.Instrument(router, config.AppName)

	if config.WithTracing {
		closeFunc := tracing.Instrument(router, config.AppName)
		config.ShutdownHooks = append(config.ShutdownHooks, closeFunc)
	}

	// add app handlers
	for _, handler := range config.Handlers {
		router.Handle(handler.Method, handler.Pattern, handler.H)
	}

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	for _, f := range config.ShutdownHooks {
		srv.RegisterOnShutdown(f)
	}

	return srv
}

// GracefulShutdown shuts down srv gracefully ane executes shutdown hooks
func GracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
