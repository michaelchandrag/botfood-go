package entities

type User struct {
	ID          int32   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Username    string  `json:"username" db:"username"`
	PhoneNumber string  `json:"phone_number" db:"phone_number"`
	Email       *string `json:"email" db:"email"`
	IsVerified  bool    `json:"is_verified" db:"is_verified"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
}
