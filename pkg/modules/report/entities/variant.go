package entities

type Variant struct {
	ID                   int     `db:"id" json:"id"`
	Slug                 string  `db:"slug" json:"slug"`
	ItemID               int     `db:"item_id" json:"item_id"`
	ItemName             string  `db:"item_name" json:"item_name"`
	VariantCategorySlug  string  `db:"variant_category_slug" json:"variant_category_slug"`
	VariantCategoryName  string  `db:"variant_category_name" json:"variant_category_name"`
	BranchChannelID      int     `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string  `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string  `db:"branch_channel_channel" json:"branch_channel_channel"`
	Name                 string  `db:"name" json:"name"`
	PayloadInStock       int     `db:"in_stock" json:"-"`
	InStock              bool    `json:"in_stock"`
	Price                *int    `db:"price" json:"price"`
	CreatedAt            *string `db:"created_at" json:"created_at"`
	UpdatedAt            *string `db:"updated_at" json:"updated_at"`
}
