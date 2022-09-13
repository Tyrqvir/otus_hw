package constraints

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/enum"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/field"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/validationerror"
)

const stringCommaSeparator = ","

func ValidateString(field *field.Field) (bool, error) {
	for _, tag := range field.Tags {
		switch tag.Name {
		case enum.LenStringKey:
			return validateStringLen(field, tag)
		case enum.RegExpStringKey:
			return validateStringRegExp(field, tag)
		case enum.InStringKey:
			return validateStringIn(field, tag)
		}
	}
	return true, nil
}

func validateStringIn(field *field.Field, tag field.Tag) (bool, error) {
	values := strings.Split(tag.Value, stringCommaSeparator)
	isValidFlag := false

	fieldVal := field.Value.String()

	for _, value := range values {
		if fieldVal == value {
			isValidFlag = true
		}
	}

	if !isValidFlag {
		return validationerror.PrepareValidationErr(validationerror.ErrNotIn, field, tag)
	}

	return true, nil
}

func validateStringRegExp(field *field.Field, tag field.Tag) (bool, error) {
	fieldValue := field.Value.String()
	_, e := regexp.MatchString(tag.RegExp.Pattern, fieldValue)
	if e != nil {
		return validationerror.PrepareValidationErr(validationerror.ErrRegExpMatching, field, tag)
	}
	return true, nil
}

func validateStringLen(field *field.Field, tag field.Tag) (bool, error) {
	valueLen := len(field.Value.String())
	ruleLen, _ := strconv.Atoi(tag.Value)

	if valueLen != ruleLen {
		return validationerror.PrepareValidationErr(validationerror.ErrStringLen, field, tag)
	}

	return true, nil
}
