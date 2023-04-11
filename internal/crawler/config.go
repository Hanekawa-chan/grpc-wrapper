package crawler

import "time"

type Config struct {
	Url     string        `envconfig:"CRAWLER_URL" default:"debug"`
	Timeout time.Duration `envconfig:"CRAWLER_TIMEOUT" default:"3m"`
}
