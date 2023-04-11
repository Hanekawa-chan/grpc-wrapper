package domain

import "context"

type Service interface {
	Search(ctx context.Context, query string) (*Company, error)
}

type Crawler interface {
	Search(ctx context.Context, query string) (*Company, error)
}

type GRPCServer interface {
	ListenAndServe() error
	Shutdown()
}
