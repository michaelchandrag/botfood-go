package dto

type OpenApiBranchChannelPromotionRequestPayload struct {
	BrandID                *int
	Keyword                string `form:"q"`
	PayloadBranchChannelID string `form:"branch_channel_id"`
	BranchChannelID        *int
	BranchChannelChannel   string `form:"branch_channel_channel"`
	BranchChannelName      string `form:"branch_channel_name"`
	DiscountType           string `form:"discount_type"`

	PayloadPage string `form:"page"`
	PayloadData string `form:"data"`
	Page        *int
	Data        *int
	SortKey     string `form:"sort_key"`
	SortValue   string `form:"sort_value"`
}
