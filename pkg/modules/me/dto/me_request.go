package dto

type MeReviewsRequestPayload struct {
	BrandID                *int
	Keyword                string   `form:"q"`
	Rating                 *int     `form:"rating"`
	QueryWithImages        string   `form:"n_images"`
	QueryWithComment       string   `form:"n_comment"`
	QueryWithMerchantReply string   `form:"n_merchant_reply"`
	InTags                 []string `form:"in_tags[]"`
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
