package dto

type ConsumerRequestPayload struct {
	MessageID   string `json:"message_id"`
	MessageSlug string
	BrandID     int                          `json:"brand_id"`
	Type        string                       `json:"type"`
	DataOutlets []ConsumerRequestDataPayload `json:"data_outlets"`
	DataItems   ConsumerDataItems            `json:"data_items"`
	RawMessage  string
}

type ConsumerRequestDataPayload struct {
	BranchChannelID      int    `json:"branch_channel_id"`
	BranchChannelName    string `json:"branch_channel_name"`
	BranchChannelChannel string `json:"branch_channel_channel"`
	PayloadIsOpen        int    `json:"is_open"`
	IsOpen               bool
	Rating               float64 `json:"rating"`
	IssuedAt             string  `json:"issued_at"`
}

type ConsumerDataItems struct {
	ItemNew    []ConsumerRequestItemDataPayload `json:"new,omitempty"`
	ItemChange []ConsumerRequestItemDataPayload `json:"change,omitempty"`
	ItemDelete []ConsumerRequestItemDataPayload `json:"delete,omitempty"`
}

type ConsumerRequestItemDataPayload struct {
	BranchChannelID      int    `json:"branch_channel_id"`
	BranchChannelName    string `json:"branch_channel_name"`
	BranchChannelChannel string `json:"branch_channel_channel"`
	ItemID               int    `json:"item_id"`
	ItemSlug             string `json:"item_slug"`
	ItemName             string `json:"item_name"`
	PayloadInStock       int    `json:"in_stock"`
	InStock              bool
	IssuedAt             string `json:"issued_at"`
}

func (dataOutlet *ConsumerRequestDataPayload) ToWebhookOutletRequestPayload() WebhookOutletRequestPayload {
	return WebhookOutletRequestPayload{
		BranchChannelID:      dataOutlet.BranchChannelID,
		BranchChannelName:    dataOutlet.BranchChannelName,
		BranchChannelChannel: dataOutlet.BranchChannelChannel,
		IsOpen:               dataOutlet.IsOpen,
		Rating:               dataOutlet.Rating,
		IssuedAt:             dataOutlet.IssuedAt,
	}
}

func (dataItem *ConsumerRequestItemDataPayload) ToWebhookItemRequestPayload() WebhookItemRequestPayload {
	return WebhookItemRequestPayload{
		BranchChannelID:      dataItem.BranchChannelID,
		BranchChannelName:    dataItem.BranchChannelName,
		BranchChannelChannel: dataItem.BranchChannelChannel,
		ItemID:               dataItem.ItemID,
		ItemSlug:             dataItem.ItemSlug,
		ItemName:             dataItem.ItemName,
		InStock:              dataItem.InStock,
		IssuedAt:             dataItem.IssuedAt,
	}
}
