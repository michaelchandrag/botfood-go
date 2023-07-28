package entities

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/michaelchandrag/botfood-go/utils"
)

type BranchChannelPromotion struct {
	ID                   int     `db:"id" json:"id"`
	Slug                 string  `db:"slug" json:"slug"`
	BranchChannelID      int     `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName    string  `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel string  `db:"branch_channel_channel" json:"branch_channel_channel"`
	Title                string  `db:"title" json:"title"`
	Description          *string `db:"description" json:"description"`
	Tags                 *string `db:"tags" json:"tags"`
	DiscountType         *string `db:"discount_type" json:"discount_type"`
	DiscountValue        *int    `db:"discount_value" json:"discount_value"`
	MinSpend             *int    `db:"min_spend" json:"min_spend"`
	MaxDiscountAmount    *int    `db:"max_discount_amount" json:"max_discount_amount"`
	CreatedAt            *string `db:"created_at" json:"created_at"`
	UpdatedAt            *string `db:"updated_at" json:"updated_at"`
}

type BranchChannelPromotionReport struct {
	ID                   int
	Slug                 string
	BranchChannelID      int
	BranchChannelName    string
	BranchChannelChannel string
	Title                string
	Description          string
	Tags                 string
	TagsInArray          []string
	TagsInText           string
	DiscountType         string
	DiscountValue        int
	MinSpend             int
	MaxDiscountAmount    int
}

func (p *BranchChannelPromotion) ToReport() (pr BranchChannelPromotionReport) {
	pr.ID = p.ID
	pr.Slug = p.Slug
	pr.BranchChannelID = p.BranchChannelID
	pr.BranchChannelName = p.BranchChannelName
	pr.BranchChannelChannel = p.BranchChannelChannel
	pr.Title = p.Title
	if p.Description != nil {
		pr.Description = *p.Description
	}
	if p.Tags != nil {
		pr.Tags = *p.Tags
		if utils.IsJSON(pr.Tags) {
			var err = json.Unmarshal([]byte(pr.Tags), &pr.TagsInArray)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				pr.TagsInText = strings.Join(pr.TagsInArray, ", ")
			}
		} else {
			pr.TagsInText = pr.Tags
		}

	}
	if p.DiscountType != nil {
		pr.DiscountType = *p.DiscountType
	}
	if p.DiscountValue != nil {
		pr.DiscountValue = *p.DiscountValue
	}
	if p.MinSpend != nil {
		pr.MinSpend = *p.MinSpend
	}
	if p.MaxDiscountAmount != nil {
		pr.MaxDiscountAmount = *p.MaxDiscountAmount
	}

	return pr
}
