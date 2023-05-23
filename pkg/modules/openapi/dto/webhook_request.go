package dto

type WebhookRequestPayload struct {
	Type        string                         `json:"type"`
	DataOutlets *[]WebhookOutletRequestPayload `json:"data_outlets,omitempty"`
	DataItems   *WebhookDataItem               `json:"data_items,omitempty"`
}

type WebhookOutletRequestPayload struct {
	BranchChannelID      int     `json:"branch_channel_id"`
	BranchChannelName    string  `json:"branch_channel_name"`
	BranchChannelChannel string  `json:"branch_channel_channel"`
	IsOpen               bool    `json:"is_open"`
	Rating               float64 `json:"rating"`
	IssuedAt             string  `json:"issued_at"`
}

type WebhookItemRequestPayload struct {
	BranchChannelID      int    `json:"branch_channel_id"`
	BranchChannelName    string `json:"branch_channel_name"`
	BranchChannelChannel string `json:"branch_channel_channel"`
	ItemSlug             string `json:"item_slug"`
	ItemID               int    `json:"item_id"`
	ItemName             string `json:"item_name"`
	InStock              bool   `json:"in_stock"`
	IssuedAt             string `json:"issued_at"`
}

type WebhookDataItem struct {
	ItemNew    *[]WebhookItemRequestPayload `json:"new,omitempty"`
	ItemChange *[]WebhookItemRequestPayload `json:"change,omitempty"`
	ItemDelete *[]WebhookItemRequestPayload `json:"delete,omitempty"`
}
