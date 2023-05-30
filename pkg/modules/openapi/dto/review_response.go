package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type OpenApiReviewListResponse struct {
	Data struct {
		Reviews     []entities.Review `json:"result"`
		CurrentPage int               `json:"current_page"`
		LimitData   int               `json:"limit_data"`
		TotalPage   int               `json:"total_page"`
		TotalData   int               `json:"total_data"`
	}
	Errors error.ErrorCollection `json:"errors"`
}
