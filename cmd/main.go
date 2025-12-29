package main

import (
	"context"
	"fmt"
	"os"

	"itkdemo/internal/repository"
	database "itkdemo/internal/transport/db"
	server "itkdemo/internal/transport/rest"
	"itkdemo/internal/usecase"
	"itkdemo/pkg/config"
	logger "itkdemo/pkg/log"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/labstack/echo/v4"
)

func gracefulShutdown(apiServer *echo.Echo, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Log.Infoln("shutting down gracefully, press Ctrl+C again to force")
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Log.Infof("Server forced to shutdown with error: %v", err)
	}

	logger.Log.Infoln("Server exiting")
	done <- true
}

func main() {
	config.Init()
	logger.Init()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	cache, err := ristretto.NewCache(&ristretto.Config[string, int64]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		logger.Log.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	db := database.New()
	repo := repository.NewPostgres(db, cache)
	service := usecase.NewWalletUseCase(repo)

	e := echo.New()
	server.NewRouter(e, service)

	go gracefulShutdown(e, done)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
	<-done
	logger.Log.Infoln("Graceful shutdown complete.")
}
