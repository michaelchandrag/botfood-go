package services

import (
	"errors"
	"fmt"
	"sync"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	item_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/item"
	variant_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/variant"
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

	variantRepository := variant_repository.NewRepository(s.db)

	var wg sync.WaitGroup
	wg.Add(len(items.Data))
	for key, item := range items.Data {
		objKey := key
		objItem := item
		go func() {
			vcs, err := variantRepository.FindByItemID(objItem.ID)
			if err != nil {
				fmt.Println(err)
			}
			items.Data[objKey].VariantCategories = vcs
			if objItem.PayloadInStock == 1 {
				items.Data[objKey].InStock = true
			} else if objItem.PayloadInStock == 0 {
				items.Data[objKey].InStock = false
			}

			wg.Done()
		}()
	}
	wg.Wait()

	response.Data.Items = items.Data
	response.Data.CurrentPage = items.CurrentPage
	response.Data.LimitData = items.LimitData
	response.Data.TotalPage = items.TotalPage
	response.Data.TotalData = items.TotalData

	return response
}
