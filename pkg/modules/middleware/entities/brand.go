package entities

type Brand struct {
	ID                 int32      `json:"id"`
	Name               string     `json:"name"`
	Slug               string     `json:"slug"`
	CreatedAt          *string    `json:"created_at"`
	UpdatedAt          *string    `json:"updated_at"`
	IsMaster           bool       `json:"is_master"`
	UserBrandThumbnail *UserBrand `json:"user_brand_thumbnail"`
}

type UserBrand struct {
	ID          int32   `json:"id"`
	UserID      int32   `json:"user_id"`
	Name        string  `json:"name"`
	Username    string  `json:"username"`
	PhoneNumber string  `json:"phone_number"`
	Email       *string `json:"email"`
	BrandID     int32   `json:"brand_id"`
	BrandSlug   string  `json:"brand_slug"`
	BrandName   string  `json:"brand_name"`
	IsThumbnail int32   `json:"is_thumbnail"`
	IsActive    int32   `json:"is_active"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}
