package entities

import "time"

type BranchChannel struct {
	ID            int                        `json:"id" db:"id"`
	BrandID       int                        `json:"-" db:"brand_id"`
	Name          string                     `json:"name" db:"name"`
	Channel       string                     `json:"channel" db:"channel"`
	PayloadIsOpen int                        `json:"-" db:"is_open"`
	IsOpen        bool                       `json:"is_open"`
	Rating        *float64                   `json:"rating" db:"rating"`
	Items         []Item                     `json:"items,omitempty"`
	Shifts        []BranchChannelShift       `json:"shifts,omitempty"`
	GroupedShifts *GroupedBranchChannelShift `json:"timing_shifts,omitempty"`
	Variants      []Variant                  `json:"variants,omitempty"`
}

type BranchChannelShift struct {
	ID              int     `json:"id" db:"id"`
	BranchChannelID int     `db:"branch_channel_id" json:"-"`
	Day             int     `json:"day" db:"day"`
	DayName         *string `json:"day_name"`
	OpenTime        string  `json:"open_time" db:"open_time"`
	CloseTime       string  `json:"close_time" db:"close_time"`
}

type GroupedBranchChannelShift struct {
	Monday    []ModernBranchChannelShift `json:"monday"`
	Tuesday   []ModernBranchChannelShift `json:"tuesday"`
	Wednesday []ModernBranchChannelShift `json:"wednesday"`
	Thursday  []ModernBranchChannelShift `json:"thursday"`
	Friday    []ModernBranchChannelShift `json:"friday"`
	Saturday  []ModernBranchChannelShift `json:"saturday"`
	Sunday    []ModernBranchChannelShift `json:"sunday"`
}

type ModernBranchChannelShift struct {
	OpenTime        string    `json:"open_time"`
	CloseTime       string    `json:"close_time"`
	OpenTimeObject  time.Time `json:"-"`
	CloseTimeObject time.Time `json:"-"`
}

func (shift *BranchChannelShift) ToModern() (modern ModernBranchChannelShift) {
	otObject, _ := time.Parse("15:04:05", shift.OpenTime)
	ctObject, _ := time.Parse("15:04:05", shift.CloseTime)
	return ModernBranchChannelShift{
		OpenTime:        shift.OpenTime,
		OpenTimeObject:  otObject,
		CloseTime:       shift.CloseTime,
		CloseTimeObject: ctObject,
	}
}
