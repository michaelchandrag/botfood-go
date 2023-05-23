package services

import (
	"errors"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	item_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/item"
)

func (s *service) GetItems(payload dto.OpenApiItemRequestPayload) (response dto.OpenApiItemListResponse) {
	itemRepository := item_repository.NewRepository(s.db)
	itemFilter := item_repository.Filter{
		BrandID:              payload.BrandID,
		Keyword:              payload.Keyword,
		InStock:              payload.InStock,
		BranchChannelID:      payload.BranchChannelID,
		BranchChannelChannel: payload.BranchChannelChannel,
		BranchChannelName:    payload.BranchChannelName,
		Name:                 payload.Name,
		Page:                 payload.Page,
		Data:                 payload.Data,
	}
	items, err := itemRepository.FindPaginated(itemFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	for key, item := range items.Data {
		if item.PayloadInStock == 1 {
			items.Data[key].InStock = true
		} else if item.PayloadInStock == 0 {
			items.Data[key].InStock = false
		}
	}

	response.Data.Items = items.Data
	response.Data.CurrentPage = items.CurrentPage
	response.Data.LimitData = items.LimitData
	response.Data.TotalPage = items.TotalPage
	response.Data.TotalData = items.TotalData

	return response
}
