package services

import (
	"errors"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	variant_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/variant"
)

func (s *service) GetVariants(payload dto.OpenApiVariantRequestPayload) (response dto.OpenApiVariantListResponse) {
	variantRepository := variant_repository.NewRepository(s.db)
	variantFilter := variant_repository.Filter{
		BrandID:              payload.BrandID,
		Keyword:              payload.Keyword,
		InStock:              payload.InStock,
		BranchChannelID:      payload.BranchChannelID,
		BranchChannelChannel: payload.BranchChannelChannel,
		BranchChannelName:    payload.BranchChannelName,
		VariantCategoryID:    payload.VariantCategoryID,
		VariantCategoryName:  payload.VariantCategoryName,
		Name:                 payload.Name,
		Page:                 payload.Page,
		Data:                 payload.Data,
	}
	variants, err := variantRepository.FindPaginated(variantFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	for key, variant := range variants.Data {
		if variant.PayloadInStock == 1 {
			variants.Data[key].InStock = true
		} else if variant.PayloadInStock == 0 {
			variants.Data[key].InStock = false
		}
	}

	response.Data.Items = variants.Data
	response.Data.CurrentPage = variants.CurrentPage
	response.Data.LimitData = variants.LimitData
	response.Data.TotalPage = variants.TotalPage
	response.Data.TotalData = variants.TotalData

	return response
}
