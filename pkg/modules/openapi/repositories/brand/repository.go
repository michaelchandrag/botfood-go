package repositories

import (
	"database/sql"
	"fmt"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindOne(filter Filter) (brand entities.Brand, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID     *int
	Slug   string
	ApiKey string
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			brands.id,
			brands.name,
			brands.slug,
			brands.webhook_url,
			brands.api_key,
			brands.secret_key,
			brands.created_at,
			brands.updated_at
		FROM
			brands
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")
	if filter.ID != nil {
		where.And("brands.id = ?", *filter.ID)
	}

	if len(filter.Slug) > 0 {
		where.And("brands.slug = ?", filter.Slug)
	}

	if len(filter.ApiKey) > 0 {
		where.And("brands.api_key = ?", filter.ApiKey)
	}

	where.And("brands.deleted_at IS NULL")

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindOne(filter Filter) (brand entities.Brand, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().QueryRowx(formattedQuery).StructScan(&brand)
	if err == sql.ErrNoRows {
		return brand, nil
	} else if err != nil {
		fmt.Println(err)
	}
	return brand, nil
}
