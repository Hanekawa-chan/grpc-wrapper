package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"rusprofwrapper/internal/crawler/models"
	"rusprofwrapper/internal/domain"
	"strings"
)

func (a *adapter) Search(_ context.Context, query string) (*domain.Company, error) {
	resp, err := a.client.Get(a.config.Url + "/ajax.php?query=" + query + "&action=search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var searchResponse models.SearchResponse
	err = json.Unmarshal(b, &searchResponse)
	if err != nil {
		return nil, err
	}

	if searchResponse.CompaniesCount < 1 || len(searchResponse.Companies) < 1 {
		return nil, errors.New("not found")
	}

	kpp, err := a.parseKpp(searchResponse.Companies[0].Link)
	if err != nil {
		return nil, err
	}

	inn := searchResponse.Companies[0].Inn[3:13]

	company := domain.Company{
		Name:    searchResponse.Companies[0].Name,
		Inn:     inn,
		Kpp:     kpp,
		CeoName: searchResponse.Companies[0].CeoName,
	}

	return &company, nil
}

func (a *adapter) parseKpp(id string) (string, error) {
	resp, err := a.client.Get(a.config.Url + id)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	clipKppIndex := strings.Index(string(b), "clip_kpp")
	tagEndIndex := clipKppIndex + strings.Index(string(b[clipKppIndex:]), ">") + 1
	if tagEndIndex+9 > len(b) {
		return "", errors.New("damaged html")
	}
	kpp := string(b[tagEndIndex : tagEndIndex+9])

	return kpp, nil
}
