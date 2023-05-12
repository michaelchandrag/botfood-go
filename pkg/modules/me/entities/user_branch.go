package entities

type UserBranch struct {
	ID         int32   `json:"id" db:"id"`
	UserID     int32   `json:"user_id" db:"user_id"`
	BranchID   int32   `json:"branch_id" db:"branch_id"`
	BranchName string  `json:"branch_name" db:"branch_name"`
	IsActive   bool    `json:"is_active" db:"is_active"`
	CreatedAt  *string `json:"created_at" db:"created_at"`
}
