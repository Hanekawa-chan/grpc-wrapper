package grpcserver

import "time"

type Config struct {
	Address           string        `envconfig:"GRPC_ADDRESS" required:"yes"`
	HttpAddress       string        `envconfig:"HTTP_ADDRESS" required:"yes"`
	MaxConnectionIdle time.Duration `envconfig:"GRPC_MAX_CONNECTION_IDLE" default:"1s"`
	Timeout           time.Duration `envconfig:"GRPC_TIMEOUT" default:"1m"`
	MaxConnectionAge  time.Duration `envconfig:"GRPC_MAX_CONNECTION_AGE" default:"1s"`
}
