package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/constraints"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/field"
	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/validationerror"
)

var (
	ValidationErrorsInstance validationerror.ValidationErrors
	mutex                    = &sync.Mutex{}
)

func Validate(v interface{}) error {
	valueOf := reflect.ValueOf(v)

	if !isStruct(valueOf) {
		return fmt.Errorf("%w %s", validationerror.ErrInvalidInputType, reflect.Struct)
	}

	if ok, e := validate(valueOf); !ok {
		return e
	}

	return nil
}

func validate(valueOf reflect.Value) (bool, error) {
	fields := getFields(valueOf)

	if len(fields) == 0 {
		return false, nil
	}

	validationErrors := make([]*validationerror.ValidationError, 0)

	for _, f := range fields {
		fmt.Printf("Validating field... [%s]\n", f.Type.Name)
		if ok, e := validateField(f); !ok {
			var validationErr *validationerror.ValidationError
			if errors.As(e, &validationErr) {
				validationErrors = append(validationErrors, validationErr)
			} else if e != nil {
				return false, e
			}
		}
	}

	if len(validationErrors) != 0 {
		mutex.Lock()
		defer mutex.Unlock()
		ValidationErrorsInstance = validationErrors
		return false, ValidationErrorsInstance
	}

	return true, nil
}

func getFields(v reflect.Value) []*field.Field {
	fields := make([]*field.Field, 0)

	for i := 0; i < v.NumField(); i++ {
		newField := field.Create(v, i)

		fieldValueKind := newField.Value.Kind()

		if fieldValueKind != reflect.Slice {
			if !newField.HasValidationTag() {
				continue
			}

			fields = append(fields, newField)
		} else {
			for j := 0; j < newField.Value.Len(); j++ {
				currentFieldName := fmt.Sprintf("%s[%d]", newField.Type.Name, j)
				currentFieldVal := newField.Value.Index(j)
				modifiedField := newField.ModifyPointerFieldForSlice(currentFieldVal, currentFieldName)
				if !newField.HasValidationTag() {
					continue
				}

				fields = append(fields, modifiedField)
			}
		}
	}

	return fields
}

func isStruct(v reflect.Value) bool {
	return reflect.Struct == v.Kind()
}

func validateField(field *field.Field) (bool, error) {
	switch field.Value.Kind() { //nolint:exhaustive
	case reflect.Int:
		return constraints.ValidateInt(field)
	case reflect.String:
		return constraints.ValidateString(field)
	}

	return false, nil
}
