package entities

import "strconv"

type Item struct {
	ID                   int     `db:"id" json:"id"`
	Slug                 string  `db:"slug" json:"slug"`
	BranchChannelID      int     `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string  `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string  `db:"branch_channel_channel" json:"branch_channel_channel"`
	Name                 string  `db:"name" json:"name"`
	Description          *string `db:"description" json:"description"`
	PayloadInStock       int     `db:"in_stock" json:"-"`
	InStock              bool    `json:"in_stock"`
	Price                *int    `db:"price" json:"price"`
	SellingPrice         *int    `db:"selling_price" json:"selling_price"`
	ImageURL             *string `db:"image_url" json:"image_url"`
	PayloadIsBundle      int     `db:"is_bundle" json:"-"`
	IsBundle             bool    `json:"is_bundle"`
	LastNotActiveAt      *string `db:"last_not_active_at" json:"last_not_active_at"`
	CreatedAt            *string `db:"created_at" json:"created_at"`
	UpdatedAt            *string `db:"updated_at" json:"updated_at"`
}

type ItemReport struct {
	ID                   int
	Slug                 string
	BranchChannelID      int
	BranchChannelName    string
	BranchChannelChannel string
	Name                 string
	Description          string
	PayloadInStock       int
	InStock              bool
	Price                int
	PriceInText          string
	SellingPrice         int
	SellingPriceInText   string
	ImageURL             string
	PayloadIsBundle      int
	IsBundle             bool
}

func (i *Item) ToReport() (ir ItemReport) {
	ir.ID = i.ID
	ir.Slug = i.Slug
	ir.BranchChannelID = i.BranchChannelID
	ir.BranchChannelName = i.BranchChannelName
	ir.BranchChannelChannel = i.BranchChannelChannel
	ir.Name = i.Name
	if i.PayloadInStock == 0 {
		ir.InStock = false
	} else {
		ir.InStock = true
	}
	if i.Price != nil {
		ir.Price = *i.Price
		ir.PriceInText = strconv.Itoa(ir.Price)
	}
	if i.SellingPrice != nil {
		ir.SellingPrice = *i.SellingPrice
		ir.SellingPriceInText = strconv.Itoa(ir.SellingPrice)
	}
	if i.ImageURL != nil {
		ir.ImageURL = *i.ImageURL
	}
	if i.Description != nil {
		ir.Description = *i.Description
	}

	return ir
}
