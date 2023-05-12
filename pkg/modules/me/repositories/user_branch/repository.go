package repositories

import (
	"fmt"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/me/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindAll(filter Filter) (userBranchs []entities.UserBranch, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	UserID   *int32
	BrandID  *int32
	IsActive bool
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			user_branchs.id,
			user_branchs.user_id,
			user_branchs.branch_id,
			branchs.name as branch_name,
			user_branchs.is_active,
			user_branchs.created_at
		FROM
			user_branchs
		JOIN branchs on branchs.id = user_branchs.branch_id AND branchs.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")
	if filter.UserID != nil {
		where.And("user_branchs.user_id = ?", *filter.UserID)
	}

	if filter.BrandID != nil {
		where.And("branchs.brand_id = ?", *filter.BrandID)
	}

	if filter.IsActive {
		where.And("user_branchs.is_active = 1")
	}

	where.And("user_branchs.deleted_at IS NULL")

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindAll(filter Filter) (userBranchs []entities.UserBranch, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&userBranchs, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return userBranchs, err
}
