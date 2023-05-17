package entities

type MessageQueue struct {
	ID        int     `json:"id" db:"id"`
	Type      string  `json:"type" db:"type"`
	BrandID   int     `json:"brand_id" db:"brand_id"`
	MessageID string  `json:"message_id" db:"message_id"`
	Body      string  `json:"body" db:"body"`
	CreatedAt *string `json:"created_at" db:"created_at"`
	UpdatedAt *string `json:"updated_at" db:"updated_at"`
	DeletedAt *string `json:"deleted_at" db:"deleted_at"`
}
