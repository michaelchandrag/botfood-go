package repositories

import (
	"fmt"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/me/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindOne(filter Filter) (user entities.User, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID         *int32
	IsVerified bool
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			users.id,
			users.name,
			users.username,
			users.phone_number,
			users.email,
			users.is_verified,
			users.created_at
		FROM
			users
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")
	if filter.ID != nil {
		where.And("users.id = ?", *filter.ID)
	}

	if filter.IsVerified {
		where.And("users.is_verified = 1")
	}

	where.And("users.deleted_at IS NULL")

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindOne(filter Filter) (user entities.User, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().QueryRowx(formattedQuery).StructScan(&user)
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}
