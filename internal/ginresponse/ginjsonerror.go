package ginresponse

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
	"net/http"
)

const DefaultHTTPErrorCode = http.StatusInternalServerError

func gatherHTTPErrorCode(myError *error) int {
	httpErrorCode := DefaultHTTPErrorCode
	if errVal, ok := (*myError).(*HTTPError); ok {
		if errVal.HTTPErrorCode != 0 {
			httpErrorCode = errVal.HTTPErrorCode
		}
	}
	return httpErrorCode
}

// GinJSONError returns a JSON API formatted HTTP error through a Gin context
func GinJSONError(ctx *gin.Context, myError error) {
	httpErrorCode := gatherHTTPErrorCode(&myError)
	errPayload := jsonapi.ErrorsPayload{Errors: []*jsonapi.ErrorObject{{Title: myError.Error()}}}
	log.Printf("%+v", errPayload)
	ctx.JSON(httpErrorCode, errPayload)
}
