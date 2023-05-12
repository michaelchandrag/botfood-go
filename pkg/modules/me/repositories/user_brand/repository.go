package repositories

import (
	"fmt"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/me/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindAll(filter Filter) (userBrands []entities.UserBrand, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	UserID   *int32
	IsActive bool
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`SELECT
			user_brands.id,
			user_brands.user_id,
			users.name,
			users.username,
			users.phone_number,
			users.email,
			user_brands.brand_id,
			brands.slug as brand_slug,
			brands.name as brand_name,
			user_brands.is_thumbnail,
			user_brands.is_active,
			user_brands.created_at
		FROM
			user_brands
		JOIN brands on brands.id = user_brands.brand_id AND brands.deleted_at IS NULL
		JOIN users on users.id = user_brands.user_id AND users.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")
	if filter.UserID != nil {
		where.And("user_brands.user_id = ?", *filter.UserID)
	}

	if filter.IsActive {
		where.And("user_brands.is_active = 1")
	}

	where.And("user_brands.deleted_at IS NULL")

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindAll(filter Filter) (userBrands []entities.UserBrand, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&userBrands, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return userBrands, err
}
