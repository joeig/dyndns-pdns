package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joeig/dyndns-pdns/internal/genericerror"
	"log"
)

func generateRequestID() (string, error) {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Unable to generate request Id")
		return "", &genericerror.GenericError{}
	}

	return uuid4.String(), nil
}

// RequestIDMiddleware adds a unique request ID to each HTTP response
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID, err := generateRequestID()
		if err != nil {
			return
		}

		ctx.Set("RequestId", requestID)
		ctx.Header("X-Request-Id", requestID)

		log.SetPrefix(fmt.Sprintf("[%s] ", requestID))
		log.Printf("Set request Id to \"%s\"", requestID)

		ctx.Next()
	}
}
