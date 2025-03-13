package genericerror

// GenericError contains information regarding a certain error
type GenericError struct {
	Message string
}

// GenericError returns an error message string
func (e *GenericError) Error() string {
	return e.Message
}
