package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	Create(mq entities.MessageQueue) (newMq entities.MessageQueue, err error)
	FindOne(filter Filter) (messageQueue entities.MessageQueue, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID          *int
	BrandID     *int
	Type        string
	MessageID   string
	MessageSlug string
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			message_queues.id,
			message_queues.brand_id,
			message_queues.type,
			message_queues.message_id,
			message_queues.body,
			message_queues.created_at,
			message_queues.updated_at,
			message_queues.deleted_at
		FROM
			message_queues
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")
	if filter.ID != nil {
		where.And("message_queues.id = ?", *filter.ID)
	}

	if filter.BrandID != nil {
		where.And("message_queues.brand_id = ?", *filter.BrandID)
	}

	if len(filter.MessageID) > 0 {
		where.And("message_queues.message_id = ?", filter.MessageID)
	}

	if len(filter.Type) > 0 {
		where.And("message_queues.type = ?", filter.Type)
	}

	where.And("message_queues.deleted_at IS NULL")

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindOne(filter Filter) (messageQueue entities.MessageQueue, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().QueryRowx(formattedQuery).StructScan(&messageQueue)
	if err == sql.ErrNoRows {
		return messageQueue, nil
	} else if err != nil {
		fmt.Println(err)
	}
	return messageQueue, nil
}

func (r *repository) Create(mq entities.MessageQueue) (newMq entities.MessageQueue, err error) {
	query := `
		INSERT INTO message_queues
			(brand_id, type, message_id, body, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	dt := time.Now()
	now := dt.Format("2006-01-02 15:04:05")
	mq.CreatedAt = &now
	mq.UpdatedAt = &now
	action := r.db.GetDB().MustExec(query, mq.BrandID, mq.Type, mq.MessageID, mq.Body, mq.CreatedAt, mq.UpdatedAt, nil)
	insertId, err := action.LastInsertId()
	newMq = mq
	newMq.ID = int(insertId)
	if err != nil {
		return newMq, err
	}

	return newMq, nil
}
