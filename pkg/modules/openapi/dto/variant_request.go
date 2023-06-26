package dto

type OpenApiVariantRequestPayload struct {
	BrandID                *int
	Keyword                string `form:"q"`
	Name                   string `form:"name"`
	PayloadBranchChannelID string `form:"branch_channel_id"`
	BranchChannelID        *int
	BranchChannelChannel   string `form:"branch_channel_channel"`
	BranchChannelName      string `form:"branch_channel_name"`
	PayloadInStock         string `form:"in_stock"`
	InStock                *int

	VariantCategoryName      string `form:"variant_category_name"`
	PayloadVariantCategoryID string `form:"variant_category_id"`
	VariantCategoryID        *int

	PayloadPage string `form:"page"`
	PayloadData string `form:"data"`
	Page        *int
	Data        *int
	SortKey     string `form:"sort_key"`
	SortValue   string `form:"sort_value"`
}
