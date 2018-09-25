package main

import "fmt"

// Error contains information regarding a certain error
type Error struct {
	Message string
}

// Error returns an error message string
func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
