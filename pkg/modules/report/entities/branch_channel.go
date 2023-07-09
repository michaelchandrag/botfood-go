package entities

const BRANCH_CHANNEL_CHANNEL_GOFOOD = "GoFood"
const BRANCH_CHANNEL_CHANNEL_GRABFOOD = "GrabFood"
const BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD = "ShopeeFood"

type BranchChannel struct {
	ID            int      `json:"id" db:"id"`
	BrandID       int      `json:"-" db:"brand_id"`
	Name          string   `json:"name" db:"name"`
	Channel       string   `json:"channel" db:"channel"`
	PayloadIsOpen int      `json:"-" db:"is_open"`
	IsOpen        bool     `json:"is_open"`
	Rating        *float64 `json:"rating" db:"rating"`
}
