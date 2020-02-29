package ginResponse

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
)

// GinJSONError returns a JSON API formatted HTTP error through a Gin context
func GinJSONError(ctx *gin.Context, httpErrorCode int, title string) {
	errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: title}}}
	log.Printf("%+v", errPayload)
	ctx.JSON(httpErrorCode, errPayload)
}
