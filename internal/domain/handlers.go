package domain

import "context"

func (s *service) Search(ctx context.Context, query string) (*Company, error) {
	return s.crawler.Search(ctx, query)
}
