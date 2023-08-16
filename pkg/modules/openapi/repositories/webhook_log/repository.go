package repositories

import (
	"time"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
)

type Repository interface {
	Create(wl entities.WebhookLog) (newWl entities.WebhookLog, err error)
}

type repository struct {
	db database.MainDB
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(wl entities.WebhookLog) (newWl entities.WebhookLog, err error) {
	query := `
		INSERT INTO webhook_logs
			(brand_id, request_url, request_body, response_body, http_response_code, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	dt := time.Now()
	now := dt.Format("2006-01-02 15:04:05")
	wl.CreatedAt = &now
	wl.UpdatedAt = &now
	action := r.db.GetDB().MustExec(query, wl.BrandID, wl.RequestURL, wl.RequestBody, wl.ResponseBody, wl.HTTPResponseCode, wl.CreatedAt, wl.UpdatedAt, nil)
	insertId, err := action.LastInsertId()
	newWl = wl
	newWl.ID = int(insertId)
	if err != nil {
		return newWl, err
	}

	return newWl, nil
}
