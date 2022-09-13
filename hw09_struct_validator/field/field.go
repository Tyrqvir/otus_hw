package field

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	validationTagKey  = "validate"
	regExpKey         = "regexp"
	tagSeparator      = "|"
	tagValueSeparator = ":"
)

type RegExp struct {
	Pattern string
	Error   error
}

type Tag struct {
	Name   string
	Value  string
	RegExp RegExp
}

type Field struct {
	Value reflect.Value
	Type  reflect.StructField
	Tags  []Tag
}

func Create(v reflect.Value, i int) *Field {
	fieldType := v.Type().Field(i)
	return &Field{
		Value: v.Field(i),
		Type:  fieldType,
		Tags:  extractTags(fieldType.Tag),
	}
}

func (field *Field) HasValidationTag() bool {
	return len(field.Tags) == 1
}

func hasValidationTag(structTag reflect.StructTag) bool {
	_, exist := structTag.Lookup(validationTagKey)
	return exist
}

func getValidationTag(structTags reflect.StructTag) string {
	return structTags.Get(validationTagKey)
}

func extractTags(structTags reflect.StructTag) []Tag {
	if !hasValidationTag(structTags) {
		return []Tag{}
	}

	tagValues := strings.Split(getValidationTag(structTags), tagSeparator)

	fieldTags := make([]Tag, 0)

	for _, tag := range tagValues {
		fieldTags = append(fieldTags, prepareTag(tag))
	}

	return fieldTags
}

func prepareTag(tag string) Tag {
	tagSchema := strings.Split(tag, tagValueSeparator)

	regExp := RegExp{
		Pattern: "",
		Error:   nil,
	}

	tagName := tagSchema[0]
	value := tagSchema[1]

	if tagName == regExpKey {
		_, err := regexp.Compile(value)
		regExp = RegExp{
			Pattern: value,
			Error:   err,
		}
	}
	return Tag{
		Name:   tagName,
		Value:  value,
		RegExp: regExp,
	}
}

func (field Field) ModifyPointerFieldForSlice(val reflect.Value, name string) *Field {
	field.Value = val
	field.Type.Name = name

	return &field
}
