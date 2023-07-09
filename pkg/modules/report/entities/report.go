package entities

type ChannelReportData struct {
	GoFoodReport ChannelReport `json:"gofood"`
}

type ChannelReport struct {
	UniqueItemNames     []string              `json:"unique_item_names"`
	BranchChannelReport []BranchChannelReport `json:"branch_channel_report,omitempty"`
}

type BranchChannelReport struct {
	IsOpen bool `json:"is_open"`
}
