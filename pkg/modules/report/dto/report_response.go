package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
	"github.com/xuri/excelize/v2"
)

type ReportResponse struct {
	Message        string                   `json:"message"`
	BranchChannels []entities.BranchChannel `json:"branch_channels"`
	Errors         error.ErrorCollection    `json:"errors"`
}

type ChannelReportResponse struct {
	Data struct {
		BranchChannels    []entities.BranchChannel   `json:"-"`
		Items             []entities.Item            `json:"-"`
		ChannelReportData entities.ChannelReportData `json:"channel_report"`
	} `json:"data"`
	File struct {
		Excel *excelize.File
	}
	Errors error.ErrorCollection `json:"errors"`
}

type PromotionReportResponse struct {
	Data struct {
		Promotions    []entities.BranchChannelPromotion `json:"promotions"`
		ItemDiscounts []entities.Item                   `json:"item_discounts"`
		ItemBundles   []entities.Item                   `json:"item_bundles"`
	} `json:"data"`
	File struct {
		Excel *excelize.File
	}
	Errors error.ErrorCollection `json:"errors"`
}

type ATPReportResponse struct {
	Data struct {
	} `json:"data"`
	File struct {
		Excel *excelize.File
	}
	Errors error.ErrorCollection `json:"errors"`
}
