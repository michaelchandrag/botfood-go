package dto

type OpenApiReportItemAvailabilityReportsRequestPayload struct {
	BrandID *int

	PayloadBranchChannelID string `form:"branch_channel_id"`
	BranchChannelID        *int
	PayloadDate            string `form:"date"`
	Date                   string

	FromCreatedAt  string `form:"from_created_at"`
	UntilCreatedAt string `form:"until_created_at"`

	Page      *int   `form:"page"`
	Data      *int   `form:"data"`
	SortKey   string `form:"sort_key"`
	SortValue string `form:"sort_value"`
}
