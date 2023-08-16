package entities

type ItemAvailabilityReport struct {
	ID                     int       `json:"id" db:"id"`
	BrandID                int       `db:"brand_id" json:"-"`
	BranchChannelID        int       `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName      string    `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel   string    `db:"branch_channel_channel" json:"branch_channel_channel"`
	ItemID                 int       `db:"item_id" json:"item_id"`
	ItemSlug               string    `db:"item_slug" json:"item_slug"`
	ItemName               string    `db:"item_name" json:"item_name"`
	Date                   string    `db:"date" json:"date"`
	ActiveTime             int       `db:"active_time" json:"active_time"`
	InactiveTime           int       `db:"inactive_time" json:"inactive_time"`
	ItemNameAtThatTime     string    `db:"item_name_at_that_time" json:"item_name_at_that_time"`
	Timeline               string    `db:"timeline" json:"timeline"`
	AvailabilityPercentage *float64  `db:"availability_percentage" json:"availability_percentage"`
	PayloadRemarks         *string   `db:"remarks" json:"-"`
	Remarks                *[]string `json:"remarks"`
	AdditionalRemarks      *string   `db:"additional_remarks" json:"additional_remarks"`
	SystemRemarks          *string   `db:"system_remarks" json:"system_remarks"`
}
