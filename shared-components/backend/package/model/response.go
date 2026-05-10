package model

type PagingResponse struct {
	TotalRecords int `json:"total_records"`
	Page         int `json:"page"`
	ItemsPerPage int `json:"items_per_page"`
}
