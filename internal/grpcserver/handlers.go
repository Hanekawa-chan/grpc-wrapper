package grpcserver

import (
	"context"
	"rusprofwrapper/protocol/services"
)

func (a *adapter) Search(ctx context.Context, request *services.SearchRequest) (*services.Company, error) {
	company, err := a.service.Search(ctx, request.Query)
	if err != nil {
		return nil, err
	}
	return &services.Company{
		Name:    company.Name,
		CeoName: company.CeoName,
		Inn:     company.Inn,
		Kpp:     company.Kpp,
	}, nil
}
