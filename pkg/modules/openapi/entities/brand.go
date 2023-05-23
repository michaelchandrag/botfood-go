package entities

type Brand struct {
	ID         int32   `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	Slug       string  `db:"slug" json:"slug"`
	WebhookURL string  `db:"webhook_url" json:"webhook_url"`
	ApiKey     *string `db:"api_key" json:"-"`
	SecretKey  *string `db:"secret_key" json:"-"`
	CreatedAt  *string `db:"created_at" json:"created_at"`
	UpdatedAt  *string `db:"updated_at" json:"updated_at"`
}
