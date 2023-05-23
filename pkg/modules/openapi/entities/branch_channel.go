package entities

type BranchChannel struct {
	ID            int      `json:"id" db:"id"`
	BrandID       int      `json:"-" db:"brand_id"`
	Name          string   `json:"name" db:"name"`
	Channel       string   `json:"channel" db:"channel"`
	PayloadIsOpen int      `json:"-" db:"is_open"`
	IsOpen        bool     `json:"is_open"`
	Rating        *float64 `json:"rating" db:"rating"`
	Items         []Item   `json:"items,omitempty"`
}
