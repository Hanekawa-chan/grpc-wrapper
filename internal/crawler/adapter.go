package crawler

import (
	"net/http"
	"rusprofwrapper/internal/domain"
)

type adapter struct {
	client *http.Client
	config *Config
}

func NewAdapter(config *Config) domain.Crawler {
	a := adapter{
		client: &http.Client{Timeout: config.Timeout},
		config: config,
	}

	return &a
}
