package entities

type WebhookLog struct {
	ID           int     `db:"id"`
	BrandID      int     `db:"brand_id"`
	RequestURL   string  `db:"request_url"`
	RequestBody  string  `db:"request_body"`
	ResponseBody string  `db:"response_body"`
	CreatedAt    *string `db:"created_at"`
	UpdatedAt    *string `db:"updated_at"`
}
