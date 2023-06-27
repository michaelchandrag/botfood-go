package entities

type Variant struct {
	ID                         int    `db:"id" json:"id"`
	VariantCategorySlug        string `db:"variant_category_slug" json:"-"`
	VariantCategoryID          int    `db:"variant_category_id" json:"variant_category_id,omitempty"`
	VariantCategoryName        string `db:"variant_category_name" json:"variant_category_name,omitempty"`
	VariantCategoryIsRequired  int    `db:"variant_category_is_required" json:"variant_category_is_required,omitempty"`
	VariantCategoryMinQuantity int    `db:"variant_category_min_quantity" json:"variant_category_min_quantity,omitempty"`
	VariantCategoryMaxQuantity int    `db:"variant_category_max_quantity" json:"variant_category_max_quantity,omitempty"`
	Slug                       string `db:"slug" json:"-"`
	BranchChannelID            int    `db:"branch_channel_id" json:"branch_channel_id,omitempty"`
	BranchChannelName          string `db:"branch_channel_name" json:"branch_channel_name,omitempty"`
	BranchChannelChannel       string `db:"branch_channel_channel" json:"branch_channel_channel,omitempty"`
	ItemVariantCategoryID      int    `db:"item_variant_category_id" json:"-"`
	ItemVariantCategoryItemID  int    `db:"item_variant_category_item_id" json:"-"`
	Name                       string `db:"name" json:"name"`
	PayloadInStock             int    `db:"in_stock" json:"-"`
	InStock                    bool   `json:"in_stock"`
	Price                      *int   `db:"price" json:"price"`
}

type VariantCategory struct {
	ID                int             `db:"id" json:"id"`
	Name              string          `db:"name" json:"name"`
	PayloadIsRequired int             `db:"is_required" json:"-"`
	IsRequired        bool            `json:"is_required"`
	MinQuantity       int             `db:"min_quantity" json:"min_quantity"`
	MaxQuantity       int             `db:"max_quantity" json:"max_quantity"`
	Variants          []ModernVariant `json:"variants"`
}

type ModernVariant struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	InStock bool   `json:"in_stock"`
	Price   *int   `db:"price" json:"price"`
}

func (v *Variant) ToModern() (modern ModernVariant) {
	inStock := false
	if v.PayloadInStock == 0 {
		inStock = false
	} else {
		inStock = true
	}
	return ModernVariant{
		ID:      v.ID,
		Name:    v.Name,
		InStock: inStock,
		Price:   v.Price,
	}
}

func (v *Variant) ToRaw() (variant Variant) {
	inStock := false
	if v.PayloadInStock == 0 {
		inStock = false
	} else {
		inStock = true
	}
	return Variant{
		ID:                  v.ID,
		VariantCategoryID:   v.VariantCategoryID,
		VariantCategoryName: v.VariantCategoryName,
		Name:                v.Name,
		InStock:             inStock,
		Price:               v.Price,
	}
}
