package entities

type UserBrand struct {
	ID          int32   `json:"id" db:"id"`
	UserID      int32   `json:"user_id" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Username    string  `json:"username" db:"username"`
	PhoneNumber string  `json:"phone_number" db:"phone_number"`
	Email       *string `json:"email" db:"email"`
	BrandID     int32   `json:"brand_id" db:"brand_id"`
	BrandSlug   string  `json:"brand_slug" db:"brand_slug"`
	BrandName   string  `json:"brand_name" db:"brand_name"`
	IsThumbnail int32   `json:"is_thumbnail" db:"is_thumbnail"`
	IsActive    int32   `json:"is_active" db:"is_active"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`
}
