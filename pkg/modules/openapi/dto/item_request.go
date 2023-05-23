package dto

type OpenApiItemRequestPayload struct {
	BrandID                *int
	Keyword                string `form:"q"`
	Name                   string `form:"name"`
	PayloadBranchChannelID string `form:"branch_channel_id"`
	BranchChannelID        *int
	BranchChannelChannel   string `form:"branch_channel_channel"`
	BranchChannelName      string `form:"branch_channel_name"`
	PayloadInStock         string `form:"in_stock"`
	InStock                *int

	Page      *int   `form:"page"`
	Data      *int   `form:"data"`
	SortKey   string `form:"sort_key"`
	SortValue string `form:"sort_value"`
}
