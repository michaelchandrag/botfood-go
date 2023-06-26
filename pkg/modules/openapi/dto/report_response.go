package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type OpenApiReportItemAvailabilityReportListResponse struct {
	Data struct {
		Report      []entities.ItemAvailabilityReport `json:"result"`
		CurrentPage int                               `json:"current_page"`
		LimitData   int                               `json:"limit_data"`
		TotalPage   int                               `json:"total_page"`
		TotalData   int                               `json:"total_data"`
	}
	Errors error.ErrorCollection `json:"errors"`
}

type OpenApiReportBranchChannelAvailabilityReportListResponse struct {
	Data struct {
		Report      []entities.BranchChannelAvailabilityReport `json:"result"`
		CurrentPage int                                        `json:"current_page"`
		LimitData   int                                        `json:"limit_data"`
		TotalPage   int                                        `json:"total_page"`
		TotalData   int                                        `json:"total_data"`
	}
	Errors error.ErrorCollection `json:"errors"`
}
