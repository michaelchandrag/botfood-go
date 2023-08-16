package entities

import "fmt"

type BranchChannelAvailabilityReport struct {
	ID                             int      `json:"id" db:"id"`
	BrandID                        int      `db:"brand_id" json:"-"`
	BranchChannelID                int      `db:"branch_channel_id" json:"branch_channel_id"`
	BranchChannelName              string   `db:"branch_channel_name" json:"branch_channel_name"`
	BranchChannelChannel           string   `db:"branch_channel_channel" json:"branch_channel_channel"`
	OperationalHours               string   `db:"operational_hours" json:"operational_hours"`
	OperationalHoursDuration       int      `db:"operational_hours_duration" json:"operational_hours_duration"`
	Date                           string   `db:"date" json:"date"`
	ActiveTime                     int      `db:"active_time" json:"active_time"`
	InactiveTime                   int      `db:"inactive_time" json:"inactive_time"`
	Timeline                       string   `db:"timeline" json:"timeline"`
	InactiveTimeline               *string  `db:"inactive_timeline" json:"inactive_timeline"`
	AverageItemsActiveTimeTotal    int      `db:"average_items_active_time_total" json:"average_items_active_time_total"`
	AverageItemsInactiveTimeTotal  int      `db:"average_items_inactive_time_total" json:"average_items_inactive_time_total"`
	OpenState                      string   `db:"open_state" json:"open_state"`
	CloseState                     string   `db:"close_state" json:"close_state"`
	PayloadHasBreak                int      `db:"has_break" json:"-"`
	HasBreak                       bool     `json:"has_break"`
	ItemAvailabilityPercentage     *float64 `db:"items_availability_percentage" json:"items_availability_percentage"`
	ItemAvailabilityPercentageText string   `json:"-"`
}

func (r *BranchChannelAvailabilityReport) ToText() {
	if r.ItemAvailabilityPercentage != nil {
		r.ItemAvailabilityPercentageText = fmt.Sprintf("%.2f", *r.ItemAvailabilityPercentage) + "%"
	}
}
