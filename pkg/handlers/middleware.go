package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/michaelchandrag/botfood-go/internal/logger"
	middlewareEntity "github.com/michaelchandrag/botfood-go/pkg/modules/middleware/entities"
	bferror "github.com/michaelchandrag/botfood-go/pkg/protocols/error"
	"github.com/michaelchandrag/botfood-go/utils"
)

type payloadClaim struct {
	Data payloadClaimData `json:"data"`
}

type payloadClaimData struct {
	Brand middlewareEntity.Brand `json:"brand"`
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

func BasicMiddleware(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {

		var customError bferror.ErrorCollection
		authorizationHeader := c.Request.Header["Authorization"]
		if len(authorizationHeader) < 1 {
			customError.AddHTTPError(401, errors.New("Unauthorized"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}
		if !strings.Contains(authorizationHeader[0], "Bearer") {
			customError.AddHTTPError(401, errors.New("Unauthorized"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authorizationHeader[0], "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			secret := utils.GetEnv("BOTFOOD_JWT_KEY", "")
			return []byte(secret), nil
		})

		if err != nil {
			logger.Agent.Info(err.Error())
			customError.AddHTTPError(401, errors.New("Unauthorized"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			jsonString, err := json.Marshal(claims)
			if err != nil {
				logger.Agent.Info(err.Error())
				customError.AddHTTPError(401, errors.New("Unauthorized"))
				h.deliverError(c, customError)
				c.Abort()
				return
			}

			payload := payloadClaim{}
			err = json.Unmarshal(jsonString, &payload)
			if err != nil {
				logger.Agent.Info(err.Error())
				customError.AddHTTPError(401, errors.New("Unauthorized"))
				h.deliverError(c, customError)
				c.Abort()
				return
			}
			c.Set("auth_brand", payload.Data.Brand)

		} else {
			logger.Agent.Info(err.Error())
			customError.AddHTTPError(401, errors.New("Unauthorized"))
			h.deliverError(c, customError)
			c.Abort()
			return
		}

		c.Next()

	}
}
