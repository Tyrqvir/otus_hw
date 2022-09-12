package constraints

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/enum"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/field"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/validationerror"
)

const IntCommaSeparator = ","

func ValidateInt(field *field.Field) (bool, error) {
	for _, tag := range field.Tags {
		switch tag.Name {
		case enum.MinIntKey:
			return validateIntMin(field, tag)
		case enum.MaxIntKey:
			return validateIntMax(field, tag)
		case enum.InIntKey:
			return validateIntIn(field, tag)
		}
	}
	return true, nil
}

func validateIntIn(field *field.Field, tag field.Tag) (bool, error) {
	values := strings.Split(tag.Value, IntCommaSeparator)

	isValidFlag := false

	fieldVal := field.Value.Int()
	for _, value := range values {
		intVal, e := strconv.ParseInt(value, 10, 32)

		if e != nil {
			return false, fmt.Errorf("%w %s", validationerror.ErrInvalidInputType, reflect.String)
		}

		if fieldVal == intVal {
			isValidFlag = true
		}
	}

	if !isValidFlag {
		return validationerror.PrepareValidationErr(validationerror.ErrNotIn, field, tag)
	}

	return true, nil
}

func validateIntMax(field *field.Field, tag field.Tag) (bool, error) {
	ruleVal, fieldVal, e := prepareMinMaxValues(field, tag)
	if e != nil {
		return false, e
	}

	if fieldVal > ruleVal {
		return validationerror.PrepareValidationErr(validationerror.ErrMaxValue, field, tag)
	}

	return true, nil
}

func prepareMinMaxValues(field *field.Field, tag field.Tag) (int64, int64, error) {
	ruleVal, e := strconv.ParseInt(tag.Value, 10, 32)
	if e != nil {
		return 0, 0, fmt.Errorf("%w %s", validationerror.ErrInvalidInputType, reflect.String)
	}

	fieldVal := field.Value.Int()
	return ruleVal, fieldVal, nil
}

func validateIntMin(field *field.Field, tag field.Tag) (bool, error) {
	ruleVal, fieldVal, e := prepareMinMaxValues(field, tag)
	if e != nil {
		return false, e
	}

	if fieldVal < ruleVal {
		return validationerror.PrepareValidationErr(validationerror.ErrMinValue, field, tag)
	}

	return true, nil
}
