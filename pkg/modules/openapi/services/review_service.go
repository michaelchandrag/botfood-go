package services

import (
	"encoding/json"
	"fmt"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/utils"

	review_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/review"
)

func (s *service) GetReviews(payload dto.OpenApiReviewsRequestPayload) (response dto.OpenApiReviewListResponse) {
	reviewRepository := review_repository.NewRepository(s.db)
	reviewFilter := review_repository.Filter{
		BrandID:              payload.BrandID,
		BranchChannelID:      payload.BranchChannelID,
		BranchChannelName:    payload.BranchChannelName,
		BranchChannelChannel: payload.BranchChannelChannel,
		Data:                 payload.Data,
		Page:                 payload.Page,
		Keyword:              payload.Keyword,
		Rating:               payload.Rating,
		WithImages:           payload.WithImages,
		WithComment:          payload.WithComment,
		WithMerchantReply:    payload.WithMerchantReply,
		FromCreatedAt:        payload.FromCreatedAt,
		UntilCreatedAt:       payload.UntilCreatedAt,
	}
	reviews, err := reviewRepository.FindPaginated(reviewFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, err)
		return response
	}

	for key, review := range reviews.Data {
		if review.RawImages != nil {
			rawImages := *review.RawImages
			if utils.IsJSON(*review.RawImages) {
				err = json.Unmarshal([]byte(rawImages), &reviews.Data[key].Images)
				if err != nil {
					fmt.Println(err)
					response.Errors.AddHTTPError(500, err)
					return response
				}
			}
		}
	}
	if reviews.Data != nil {
		response.Data.Reviews = reviews.Data
	} else {
		response.Data.Reviews = []entities.Review{}
	}
	response.Data.CurrentPage = reviews.CurrentPage
	response.Data.LimitData = reviews.LimitData
	response.Data.TotalPage = reviews.TotalPage
	response.Data.TotalData = reviews.TotalData
	return response
}
