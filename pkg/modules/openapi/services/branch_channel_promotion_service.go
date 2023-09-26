package services

import (
	"encoding/json"
	"errors"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	branch_channel_promotion_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/branch_channel_promotion"
)

func (s *service) GetBranchChannelPromotions(payload dto.OpenApiBranchChannelPromotionRequestPayload) (response dto.OpenApiBranchChannelPromotionListResponse) {
	promotionRepository := branch_channel_promotion_repository.NewRepository(s.db)
	promotionFilter := branch_channel_promotion_repository.Filter{
		BrandID:              payload.BrandID,
		Keyword:              payload.Keyword,
		BranchChannelID:      payload.BranchChannelID,
		BranchChannelChannel: payload.BranchChannelChannel,
		BranchChannelName:    payload.BranchChannelName,
		DiscountType:         payload.DiscountType,
		Page:                 payload.Page,
		Data:                 payload.Data,
	}
	promotions, err := promotionRepository.FindPaginated(promotionFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	for key, promotion := range promotions.Data {
		if promotion.PayloadTags != nil {
			byteString := *promotion.PayloadTags
			err = json.Unmarshal([]byte(byteString), &promotions.Data[key].Tags)
			if err != nil {
				response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
				return response
			}
		}
	}

	response.Data.Promotions = promotions.Data
	response.Data.CurrentPage = promotions.CurrentPage
	response.Data.LimitData = promotions.LimitData
	response.Data.TotalPage = promotions.TotalPage
	response.Data.TotalData = promotions.TotalData

	return response
}
