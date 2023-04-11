package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"rusprofwrapper/internal/crawler/models"
	"rusprofwrapper/internal/domain"
	"strconv"
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
	fmt.Println(string(b))

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

	inn, err := strconv.Atoi(searchResponse.Companies[0].Inn[3:13])
	if err != nil {
		return nil, err
	}

	company := domain.Company{
		Name:    searchResponse.Companies[0].Name,
		Inn:     int32(inn),
		Kpp:     kpp,
		CeoName: searchResponse.Companies[0].CeoName,
	}

	return &company, nil
}

func (a *adapter) parseKpp(id string) (int32, error) {
	resp, err := a.client.Get(a.config.Url + id)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	clipKppIndex := strings.Index(string(b), "clip_kpp")
	tagEndIndex := clipKppIndex + strings.Index(string(b[clipKppIndex:]), ">") + 1
	if tagEndIndex+9 > len(b) {
		return 0, errors.New("damaged html")
	}
	kpp := string(b[tagEndIndex : tagEndIndex+9])

	kppInt, err := strconv.Atoi(kpp)
	if err != nil {
		return 0, err
	}

	return int32(kppInt), nil
}
