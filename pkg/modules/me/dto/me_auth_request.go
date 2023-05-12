package dto

import (
	middlewareEntity "github.com/michaelchandrag/botfood-go/pkg/modules/middleware/entities"
)

type MeAuthRequestPayload struct {
	AuthBrand middlewareEntity.Brand
}
