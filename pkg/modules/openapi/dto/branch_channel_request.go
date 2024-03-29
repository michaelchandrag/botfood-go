package dto

type OpenApiBranchChannelRequestPayload struct {
	ID            *int `form:"id"`
	BrandID       *int
	Keyword       string `form:"q"`
	Name          string `form:"name"`
	Channel       string `form:"channel"`
	PayloadIsOpen string `form:"is_open"`
	IsOpen        *int

	Page        *int
	Data        *int
	PayloadPage string `form:"page"`
	PayloadData string `form:"data"`
	SortKey     string `form:"sort_key"`
	SortValue   string `form:"sort_value"`
}
