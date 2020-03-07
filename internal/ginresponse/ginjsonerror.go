package ginresponse

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
	"net/http"
)

// GinJSONError returns a JSON API formatted HTTP error through a Gin context
func GinJSONError(ctx *gin.Context, myError error) {
	httpErrorCode := http.StatusInternalServerError
	if errVal := myError.(*HTTPError); errVal.HTTPErrorCode != 0 {
		httpErrorCode = errVal.HTTPErrorCode
	}

	errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: myError.Error()}}}
	log.Printf("%+v", errPayload)
	ctx.JSON(httpErrorCode, errPayload)
}
