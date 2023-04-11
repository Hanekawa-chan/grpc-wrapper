package grpcserver

import (
	"context"
	"github.com/Hanekawa-chan/grpc-server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"rusprofwrapper/internal/domain"
	"rusprofwrapper/protocol/services"
	"time"
)

type adapter struct {
	config  *Config
	server  *grpc.Server
	service domain.Service
}

func NewAdapter(config *Config, service domain.Service) domain.GRPCServer {
	a := &adapter{
		config:  config,
		service: service,
	}

	server := newServer(config)
	a.server = server

	return a
}

func (a *adapter) ListenAndServe() error {
	errChan := make(chan error)
	services.RegisterWrapperServiceServer(a.server, a)

	lis, err := net.Listen("tcp", a.config.Address)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen")
		return err
	}

	go func() {
		log.Log().Msg("public grpc server started on " + a.config.Address)
		if err := a.server.Serve(lis); err != nil {
			log.Error().Err(err).Msg("listen server")
			errChan <- err
		}
	}()

	time.Sleep(1 * time.Second)
	go func() {
		ctx := context.Background()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err := services.RegisterWrapperServiceHandlerFromEndpoint(ctx, mux, a.config.Address, opts)
		if err != nil {
			log.Error().Err(err).Msg("can't register http handler")
			errChan <- err
		}

		log.Log().Msg("public http server started on " + a.config.HttpAddress)
		if err := http.ListenAndServe(a.config.HttpAddress, mux); err != nil {
			panic(err)
		}
	}()

	return <-errChan
}

func (a *adapter) Shutdown() {
	a.server.GracefulStop()
	log.Info().Msg("Server Exited Properly")
}

func newServer(cfg *Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	config := &grpc_server.Config{
		MaxConnectionIdle: cfg.MaxConnectionIdle,
		Timeout:           cfg.Timeout,
		MaxConnectionAge:  cfg.MaxConnectionAge,
	}

	return grpc_server.NewGRPCServer(config, middlewares...)
}
