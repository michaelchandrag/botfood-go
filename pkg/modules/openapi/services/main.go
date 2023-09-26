package services

import (
	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
)

type OpenApiService interface {
	ConsumeActivityMessage(payload dto.ConsumerRequestPayload) (response dto.ConsumerResponsePayload)
	GetBranchChannels(payload dto.OpenApiBranchChannelRequestPayload) (response dto.OpenApiBranchChannelListResponse)
	GetItems(payload dto.OpenApiItemRequestPayload) (response dto.OpenApiItemListResponse)
	GetBranchChannelDetail(payload dto.OpenApiBranchChannelRequestPayload) (response dto.OpenApiBranchChannelDetailResponse)
	GetReviews(payload dto.OpenApiReviewsRequestPayload) (response dto.OpenApiReviewListResponse)
	GetItemAvailabilityReports(payload dto.OpenApiReportItemAvailabilityReportsRequestPayload) (response dto.OpenApiReportItemAvailabilityReportListResponse)
	GetBranchChannelAvailabilityReports(payload dto.OpenApiReportBranchChannelAvailabilityReportsRequestPayload) (response dto.OpenApiReportBranchChannelAvailabilityReportListResponse)
	GetVariants(payload dto.OpenApiVariantRequestPayload) (response dto.OpenApiVariantListResponse)
	GetBranchChannelPromotions(payload dto.OpenApiBranchChannelPromotionRequestPayload) (response dto.OpenApiBranchChannelPromotionListResponse)
}

type service struct {
	db database.MainDB
}

func RegisterOpenApiService(db database.MainDB) OpenApiService {
	return &service{
		db: db,
	}
}
