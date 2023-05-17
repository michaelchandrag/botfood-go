package dto

type ConsumerRequestPayload struct {
	MessageID   string `json:"message_id"`
	MessageSlug string
	BrandID     int                          `json:"brand_id"`
	Type        string                       `json:"type"`
	DataOutlet  []ConsumerRequestDataPayload `json:"data_outlet"`
	DataItem    struct {
		ItemNew    []ConsumerRequestItemDataPayload `json:"new"`
		ItemChange []ConsumerRequestItemDataPayload `json:"change"`
		ItemDelete []ConsumerRequestItemDataPayload `json:"delete"`
	} `json:"data_item"`
	RawMessage string
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

type ConsumerRequestItemDataPayload struct {
	BranchChannelID      int    `json:"branch_channel_id"`
	BranchChannelName    string `json:"branch_channel_name"`
	BranchChannelChannel string `json:"branch_channel_channel"`
	ItemSlug             string `json:"item_slug"`
	PayloadInStock       int    `json:"in_stock"`
	InStock              bool
	IssuedAt             string `json:"issued_at"`
}
