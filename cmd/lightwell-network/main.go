package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lightwell-network/lightwell-network-backend/pkg/config"
	"github.com/lightwell-network/lightwell-network-backend/pkg/router"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Load()
	config.ConfigureLogging()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-exit
		cancel()
	}()

	apiServer(ctx, &wg)
	wg.Wait()
}

func apiServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(2) // api server & shutdown monitor

	e := router.ConfigureEcho()

	go func() {
		defer wg.Done()
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
		log.Info().Msg("apiServer stopped")
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info().Msg("Caught context done, closing api server.")
		if err := e.Shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Fatal().Err(err).Msg("Failed to shut down server")
		}
	}()
}
