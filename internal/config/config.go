package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"rusprofwrapper/internal/crawler"
	"rusprofwrapper/internal/grpcserver"
)

type Config struct {
	Crawler    *crawler.Config
	GRPCServer *grpcserver.Config
}

func Parse() (*Config, error) {
	cfg := Config{}
	grpc := grpcserver.Config{}
	crawlerCfg := crawler.Config{}
	project := "WRAPPER"

	err := envconfig.Process(project, &grpc)
	if err != nil {
		log.Err(err).Msg("logger config error")
		return nil, err
	}
	err = envconfig.Process(project, &crawlerCfg)
	if err != nil {
		log.Err(err).Msg("logger config error")
		return nil, err
	}

	cfg.GRPCServer = &grpc
	cfg.Crawler = &crawlerCfg
	return &cfg, nil
}
