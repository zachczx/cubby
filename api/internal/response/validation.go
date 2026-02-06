package response

import "fmt"

const (
	MaxCharLength  = 255
	LongTextLength = 1028
)

type ValidationError struct {
	Field   string
	Message string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

func ValErr(field, msg string) error {
	return &ValidationError{Field: field, Message: msg}
}

func ValErrf(field string, format string, args ...any) error {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf(format, args...),
	}
}
