package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/config"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/instrumentation"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Load()
	config.ConfigureLogging()

	reg := prometheus.NewRegistry()
	metrics := instrumentation.NewMetrics(reg)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-exit
		cancel()
	}()

	apiServer(ctx, &wg, metrics)
	metricsServer(ctx, &wg, metrics)
	wg.Wait()
}

func apiServer(ctx context.Context, wg *sync.WaitGroup, metrics *instrumentation.Metrics) {
	wg.Add(2) // api server & shutdown monitor

	e := router.ConfigureEchoWithMetrics(metrics)

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

func metricsServer(ctx context.Context, wg *sync.WaitGroup, metrics *instrumentation.Metrics) {
	wg.Add(2) // metrics server & shutdown monitor

	cfg := config.Get().Metrics
	e := echo.New()
	e.HideBanner = true
	e.Add(http.MethodGet, cfg.Path, echo.WrapHandler(promhttp.HandlerFor(
		metrics.Registry(),
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Registry:          metrics.Registry(),
		},
	)))

	addr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		defer wg.Done()
		log.Info().Msgf("Starting metrics server on %s%s", addr, cfg.Path)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start metrics server")
		}
		log.Info().Msg("metricsServer stopped")
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info().Msg("Caught context done, closing metrics server.")
		if err := e.Shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Fatal().Err(err).Msg("Failed to shut down metrics server")
		}
	}()
}
