package models

type SearchResponse struct {
	CompaniesCount int `json:"ul_count"`
	Companies      []struct {
		Name    string `json:"name"`
		Link    string `json:"link"`
		Inn     string `json:"inn"`
		CeoName string `json:"ceo_name"`
	} `json:"ul"`
}
