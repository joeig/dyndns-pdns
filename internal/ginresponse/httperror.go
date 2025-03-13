package ginresponse

// HTTPError contains information regarding a certain error
type HTTPError struct {
	Message       string
	HTTPErrorCode int
}

// HTTPError returns an error message string
func (e *HTTPError) Error() string {
	return e.Message
}
