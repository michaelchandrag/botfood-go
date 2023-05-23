package handlers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	brand_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/brand"
	bferror "github.com/michaelchandrag/botfood-go/pkg/protocols/error"
	"github.com/michaelchandrag/botfood-go/utils"
)

func OpenApiMiddleware(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {

		var customError bferror.ErrorCollection

		tokenPartner := getHeader(c, "X-Partner-Token")
		if tokenPartner == "" {
			customError.AddHTTPError(401, errors.New("Unauthorized. X-Partner-Token is required"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		apiKey := getHeader(c, "X-Api-Key")
		if apiKey == "" {
			customError.AddHTTPError(401, errors.New("Unauthorized. X-Api-Key is required"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		timestamp := getHeader(c, "X-Timestamp")
		if timestamp == "" {
			customError.AddHTTPError(401, errors.New("Unauthorized. X-Timestamp is required"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		brandRepository := brand_repository.NewRepository(h.dao)
		brandFilter := brand_repository.Filter{
			ApiKey: apiKey,
		}
		brand, err := brandRepository.FindOne(brandFilter)
		if err != nil {
			customError.AddHTTPError(401, errors.New("Unauthorized. X-Api-Key is required"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		recipe := fmt.Sprintf("%s:%s", apiKey, timestamp)
		expectedToken := utils.GenerateHMAC256(recipe, *brand.SecretKey)

		if expectedToken != tokenPartner {
			fmt.Println("Expected Token: ", expectedToken)
			fmt.Println("Given Token: ", tokenPartner)
			customError.AddHTTPError(401, errors.New("Unauthorized. X-Partner-Token invalid"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		c.Set("open_api_brand", brand)

		c.Next()

	}
}

func getHeader(c *gin.Context, key string) string {
	headers := c.Request.Header[key]
	if len(headers) > 0 {
		return headers[0]
	}
	return ""
}
