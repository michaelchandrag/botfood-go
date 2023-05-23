package services

import (
	"errors"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	branch_channel_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/branch_channel"
	item_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/item"
)

func (s *service) GetBranchChannels(payload dto.OpenApiBranchChannelRequestPayload) (response dto.OpenApiBranchChannelListResponse) {
	branchChannelRepository := branch_channel_repository.NewRepository(s.db)
	branchChannelFilter := branch_channel_repository.Filter{
		BrandID: payload.BrandID,
		Keyword: payload.Keyword,
		IsOpen:  payload.IsOpen,
		Channel: payload.Channel,
		Name:    payload.Name,
		Data:    payload.Data,
		Page:    payload.Page,
	}
	branchChannels, err := branchChannelRepository.FindPaginated(branchChannelFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	for key, branchChannel := range branchChannels.Data {
		if branchChannel.PayloadIsOpen == 1 {
			branchChannels.Data[key].IsOpen = true
		} else if branchChannel.PayloadIsOpen == 0 {
			branchChannels.Data[key].IsOpen = false
		}
	}

	response.Data.BranchChannels = branchChannels.Data
	response.Data.CurrentPage = branchChannels.CurrentPage
	response.Data.LimitData = branchChannels.LimitData
	response.Data.TotalPage = branchChannels.TotalPage
	response.Data.TotalData = branchChannels.TotalData

	return response
}

func (s *service) GetBranchChannelDetail(payload dto.OpenApiBranchChannelRequestPayload) (response dto.OpenApiBranchChannelDetailResponse) {
	branchChannelRepository := branch_channel_repository.NewRepository(s.db)
	branchChannelFilter := branch_channel_repository.Filter{
		BrandID: payload.BrandID,
		ID:      payload.ID,
	}
	branchChannel, err := branchChannelRepository.FindOne(branchChannelFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	if branchChannel.ID == 0 {
		response.Errors.AddHTTPError(400, errors.New("Branch channel is required"))
		return response
	}

	if branchChannel.PayloadIsOpen == 1 {
		branchChannel.IsOpen = true
	} else if branchChannel.PayloadIsOpen == 0 {
		branchChannel.IsOpen = false
	}

	itemRepository := item_repository.NewRepository(s.db)
	itemFilter := item_repository.Filter{
		BranchChannelID: &branchChannel.ID,
	}
	items, err := itemRepository.FindAll(itemFilter)
	for key, item := range items {
		if item.PayloadInStock == 1 {
			items[key].InStock = true
		} else if item.PayloadInStock == 0 {
			items[key].InStock = false
		}
	}
	branchChannel.Items = items

	response.Data = branchChannel

	return response
}
