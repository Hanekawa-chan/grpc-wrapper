package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"rusprofwrapper/internal/config"
	"rusprofwrapper/internal/crawler"
	"rusprofwrapper/internal/domain"
	"rusprofwrapper/internal/grpcserver"
	"syscall"
)

func main() {
	// Parse all configs form env
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal().Err(err)
	}

	cr := crawler.NewAdapter(cfg.Crawler)

	service := domain.NewService(cr)
	grpcServer := grpcserver.NewAdapter(cfg.GRPCServer, service)
	log.Info().Msg("Initialized everything")

	// Channels for errors and os signals
	stop := make(chan error, 1)
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	// Receive errors form start bot func into error channel
	go func(stop chan<- error) {
		stop <- grpcServer.ListenAndServe()
	}(stop)

	// Blocking select
	select {
	case sig := <-osSig:
		log.Info().Msgf("Received os syscall signal %v", sig)
	case err := <-stop:
		log.Error().Err(err).Msg("Received Error signal")
	}

	// Shutdown code
	log.Info().Msg("Shutting down...")

	grpcServer.Shutdown()

	log.Info().Msg("Shutdown - success")
}
