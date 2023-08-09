package lib

import (
	"fmt"
	"path/filepath"
	"reflect"
)

func GetRelPath(absPath string) (*string, error) {
	base, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("Error to get an absolute path")
	}

	rel, err := filepath.Rel(base, absPath)
	if err != nil {
		return nil, fmt.Errorf("Error to get an relative path")
	}

	return &rel, nil
}

func IsEqualWithoutDescriptions(a, b interface{}) bool {
	valA, valB := reflect.ValueOf(a), reflect.ValueOf(b)

	// Check if either value is a zero Value or is not valid.
	if !valA.IsValid() || !valB.IsValid() {
		return false
	}

	// Ensure the two values are of the same type.
	if valA.Type() != valB.Type() {
		return false
	}

	// Dereference pointers to their underlying values.
	if valA.Kind() == reflect.Ptr && valA.Elem().IsValid() {
		valA = valA.Elem()
	}

	if valB.Kind() == reflect.Ptr && valB.Elem().IsValid() {
		valB = valB.Elem()
	}

	switch valA.Kind() {
	// Handle slice separately.
	case reflect.Slice:
		if valA.Len() != valB.Len() {
			return false
		}
		for i := 0; i < valA.Len(); i++ {
			if !IsEqualWithoutDescriptions(valA.Index(i).Interface(), valB.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Struct:
		typ := valA.Type()
		for i := 0; i < valA.NumField(); i++ {
			field := typ.Field(i)

			// Ignore Descriptions and Comments field.
			if field.Name == "Descriptions" || field.Name == "Comments" || field.Name == "BaseFileInfo" {
				continue
			}

			fieldA, fieldB := valA.Field(i), valB.Field(i)

			switch field.Type.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Struct:
				if !IsEqualWithoutDescriptions(fieldA.Interface(), fieldB.Interface()) {
					return false
				}
			default:
				if !fieldA.IsValid() || !fieldB.IsValid() {
					return false
				}

				if fieldA.Interface() != fieldB.Interface() {
					return false
				}
			}
		}
		return true
	default:
		return valA.Interface() == valB.Interface()
	}
}

func mergeStrings(a, b *[]string) *[]string {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	result := append(*a, *b...)
	return &result
}

func mergeDescriptionsAndComments(a, b interface{}) {
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	// Check if a and b are pointers or interfaces before calling Elem
	if !valA.IsValid() || (valA.Kind() != reflect.Ptr && valA.Kind() != reflect.Interface) {
		return
	}
	if !valB.IsValid() || (valB.Kind() != reflect.Ptr && valB.Kind() != reflect.Interface) {
		return
	}

	valA = valA.Elem()
	valB = valB.Elem()

	typ := valA.Type()

	for i := 0; i < valA.NumField(); i++ {
		field := typ.Field(i)

		if field.Name == "Descriptions" || field.Name == "Comments" {
			aField, okA := valA.Field(i).Interface().(*[]string)
			bField, okB := valB.Field(i).Interface().(*[]string)
			if okA && okB {
				valA.Field(i).Set(reflect.ValueOf(mergeStrings(aField, bField)))
			}
			continue
		}

		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr {
			for j := 0; j < valA.Field(i).Len(); j++ {
				mergeDescriptionsAndComments(valA.Field(i).Index(j).Interface(), valB.Field(i).Index(j).Interface())
			}
		}
	}
}
