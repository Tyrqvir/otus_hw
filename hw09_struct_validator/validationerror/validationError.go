package validationerror

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/field"
)

var (
	ErrInvalidInputType = errors.New("input type is not")
	ErrNotIn            = errors.New("value not in")
	ErrMinValue         = errors.New("value must be gte")
	ErrMaxValue         = errors.New("value must be lte")
	ErrRegExpMatching   = errors.New("value must be has valid pattern")
	ErrStringLen        = errors.New("string must be equal")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	var builder strings.Builder
	builder.WriteString("field: ")
	builder.WriteString(fmt.Sprintf("[%s]", v.Field))
	builder.WriteString(" - ")
	builder.WriteString(v.Err.Error())

	return builder.String()
}

type ValidationErrors []*ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, e := range v {
		builder.WriteString("field: ")
		builder.WriteString(fmt.Sprintf("[%s]", e.Field))
		builder.WriteString(" - ")
		builder.WriteString(e.Err.Error())
	}
	return builder.String()
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}

func PrepareValidationErr(e error, field *field.Field, tag field.Tag) (bool, error) {
	return false, NewValidationError(
		field.Type.Name,
		fmt.Errorf("%w : [%v], taken - [%v]", e, tag.Value, field.Value.Interface()))
}
