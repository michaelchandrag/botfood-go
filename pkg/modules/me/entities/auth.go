package entities

type Auth struct {
	Brand              Brand         `json:"brand"`
	IsMaster           bool          `json:"is_master"`
	UserBrandThumbnail *UserBrand    `json:"user_brand_thumbnail"`
	UserBrands         *[]UserBrand  `json:"user_brands"`
	UserBranchs        *[]UserBranch `json:"user_branchs"`
	User               *User         `json:"user"`
	BranchIDs          []int
}
