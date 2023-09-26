package entities

type BranchChannelPromotion struct {
	ID                   int      `db:"id" json:"id"`
	Slug                 string   `db:"slug" json:"slug"`
	BranchChannelID      int      `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string   `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string   `db:"branch_channel_channel" json:"branch_channel_channel"`
	Title                string   `db:"title" json:"title"`
	Description          *string  `db:"description" json:"description"`
	PayloadTags          *string  `db:"tags" json:"-"`
	Tags                 []string `json:"tags"`
	DiscountType         *string  `db:"discount_type" json:"discount_type"`
	DiscountValue        *int     `db:"discount_value" json:"discount_value"`
	MinSpend             *int     `db:"min_spend" json:"min_spend"`
	MaxDiscountAmount    *int     `db:"max_discount_amount" json:"max_discount_amount"`
	CreatedAt            *string  `db:"created_at" json:"created_at"`
	UpdatedAt            *string  `db:"updated_at" json:"updated_at"`
}
