package domain

type service struct {
	crawler Crawler
}

func NewService(crawler Crawler) Service {
	return &service{
		crawler: crawler,
	}
}
