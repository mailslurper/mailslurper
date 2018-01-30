package mailslurper

import "fmt"

/*
An MissingHeaderError is used to tell a caller that a requested
header has not been captured or does not exist
*/
type MissingHeaderError struct {
	MissingHeader string
}

/*
MissingHeader returns a new error object
*/
func MissingHeader(header string) *MissingHeaderError {
	return &MissingHeaderError{
		MissingHeader: header,
	}
}

func (err *MissingHeaderError) Error() string {
	return fmt.Sprintf("Missing header named '%s'", err.MissingHeader)
}
