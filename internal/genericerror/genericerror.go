package genericerror

import (
	"fmt"
)

// GenericError contains information regarding a certain error
type GenericError struct {
	Message string
}

// GenericError returns an error message string
func (e *GenericError) Error() string {
	return fmt.Sprintf(e.Message)
}
