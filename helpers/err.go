package helpers

import "fmt"

type HttpError struct {
	Code    int
	Message string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("Error Code: %d, Message: %s", e.Code, e.Message)
}

func (e HttpError) GetFields() (int, string) {
	return e.Code, e.Message
}
