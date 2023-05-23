package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type OpenApiBranchChannelListResponse struct {
	Data struct {
		BranchChannels []entities.BranchChannel `json:"result"`
		CurrentPage    int                      `json:"current_page"`
		LimitData      int                      `json:"limit_data"`
		TotalPage      int                      `json:"total_page"`
		TotalData      int                      `json:"total_data"`
	} `json:"data"`
	Errors error.ErrorCollection `json:"errors"`
}

type OpenApiBranchChannelDetailResponse struct {
	Data   entities.BranchChannel `json:"data"`
	Errors error.ErrorCollection  `json:"errors"`
}
