package ginresponse

import (
	"fmt"
)

// HTTPError contains information regarding a certain error
type HTTPError struct {
	Message       string
	HTTPErrorCode int
}

// HTTPError returns an error message string
func (e *HTTPError) Error() string {
	return fmt.Sprintf(e.Message)
}
