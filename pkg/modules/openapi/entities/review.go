package entities

type Review struct {
	ID   int    `db:"id" json:"id"`
	Slug string `db:"slug" json:"slug"`

	BranchChannelID      int    `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string `db:"branch_channel_channel" json:"branch_channel_channel"`
	BrandID              int    `db:"brand_id" json:"-"`

	Rating          int       `db:"rating" json:"rating"`
	Comment         *string   `db:"comment" json:"comment"`
	ItemName        *string   `db:"item_name" json:"item_name"`
	RawImages       *string   `db:"images" json:"-"`
	Images          *[]string `json:"images"`
	MerchantReply   *string   `db:"merchant_reply" json:"merchant_reply"`
	MerchantReplyAt *string   `db:"merchant_reply_at" json:"merchant_reply_at"`
	IsReviewed      int       `db:"is_reviewed" json:"-"`
	OrderedAt       *string   `db:"ordered_at" json:"ordered_at"`
	CreatedAt       *string   `db:"created_at" json:"created_at"`
	UpdatedAt       *string   `db:"updated_at" json:"-"`
	DeletedAt       *string   `db:"deleted_at" json:"-"`
}
