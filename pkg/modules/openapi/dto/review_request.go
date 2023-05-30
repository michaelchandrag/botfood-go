package dto

type OpenApiReviewsRequestPayload struct {
	BrandID *int

	PayloadBranchChannelID string `form:"branch_channel_id"`
	BranchChannelID        *int
	BranchChannelName      string `form:"branch_channel_name"`
	BranchChannelChannel   string `form:"branch_channel_channel"`
	Keyword                string `form:"q"`
	PayloadRating          string `form:"rating"`
	Rating                 *int
	QueryWithImages        string `form:"n_images"`
	QueryWithComment       string `form:"n_comment"`
	QueryWithMerchantReply string `form:"n_merchant_reply"`
	WithImages             bool
	WithComment            bool
	WithMerchantReply      bool
	BranchIDs              []int

	FromCreatedAt  string `form:"from_created_at"`
	UntilCreatedAt string `form:"until_created_at"`

	Page      *int   `form:"page"`
	Data      *int   `form:"data"`
	SortKey   string `form:"sort_key"`
	SortValue string `form:"sort_value"`
}
