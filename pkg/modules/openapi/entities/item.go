package entities

type Item struct {
	ID                   int               `db:"id" json:"id"`
	Slug                 string            `db:"slug" json:"slug"`
	BranchChannelID      int               `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string            `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string            `db:"branch_channel_channel" json:"branch_channel_channel"`
	Name                 string            `db:"name" json:"name"`
	PayloadInStock       int               `db:"in_stock" json:"-"`
	InStock              bool              `json:"in_stock"`
	Price                *int              `db:"price" json:"price"`
	SellingPrice         *int              `db:"selling_price" json:"selling_price"`
	ImageURL             *string           `db:"image_url" json:"image_url"`
	LastNotActiveAt      *string           `db:"last_not_active_at" json:"last_not_active_at"`
	VariantCategories    []VariantCategory `json:"variant_categories"`
	CreatedAt            *string           `db:"created_at" json:"created_at"`
	UpdatedAt            *string           `db:"updated_at" json:"updated_at"`
}
