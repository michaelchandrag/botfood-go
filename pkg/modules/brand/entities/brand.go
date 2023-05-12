package entities

type Brand struct {
	ID        int32   `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Slug      string  `db:"slug" json:"slug"`
	CreatedAt *string `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
}
