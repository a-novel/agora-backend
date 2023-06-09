package validation

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"reflect"
	"regexp"
)

func CheckRequire(field string, value interface{}) error {
	r := reflect.ValueOf(value)
	switch r.Kind() {
	case reflect.Pointer:
		if value == nil {
			return NewErrNil(field)
		}

		return CheckRequire(field, r.Elem().Interface())
	case reflect.Map, reflect.Array, reflect.Slice, reflect.String:
		if r.Len() == 0 {
			return NewErrNil(field)
		}
	default:
		if value == nil {
			return NewErrNil(field)
		}
	}

	return nil
}

func CheckMinMax(field string, value interface{}, min int, max int) error {
	r := reflect.ValueOf(value)
	switch r.Kind() {
	case reflect.Pointer:
		if value != nil {
			return CheckMinMax(field, r.Elem().Interface(), min, max)
		}
	case reflect.String:
		// Default len() methods do not give visually accurate character counts for non-ASCII strings.
		// https://stackoverflow.com/a/12668840/9021186
		var (
			ia norm.Iter
			nc int
		)
		ia.InitString(norm.NFKD, r.String())
		for !ia.Done() {
			nc = nc + 1
			ia.Next()
		}

		if min >= 0 && nc < min {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot contain less than %v characters (currently has %v)", min, nc),
			)
		}

		if max >= 0 && nc > max {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot contain more than %v characters (currently has %v)", max, nc),
			)
		}
	case reflect.Map, reflect.Array, reflect.Slice:
		if min >= 0 && r.Len() < min {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot contain less than %v elements (currently has %v)", min, r.Len()),
			)
		}

		if max >= 0 && r.Len() > max {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot contain more than %v elements (currently has %v)", max, r.Len()),
			)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if min >= 0 && r.Int() < int64(min) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be greater than %v (currently %v)", min, r.Int()),
			)
		}

		if max >= 0 && r.Int() > int64(max) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be less than %v (currently %v)", min, r.Int()),
			)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if min >= 0 && r.Uint() < uint64(min) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be greater than %v (currently %v)", min, r.Uint()),
			)
		}

		if max >= 0 && r.Uint() > uint64(max) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be less than %v (currently %v)", min, r.Uint()),
			)
		}
	case reflect.Float32, reflect.Float64:
		if min >= 0 && r.Float() < float64(min) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be greater than %v (currently %v)", min, r.Float()),
			)
		}

		if max >= 0 && r.Float() > float64(max) {
			return NewErrInvalidEntity(
				field,
				fmt.Sprintf("it cannot be less than %v (currently %v)", min, r.Float()),
			)
		}
	}

	return nil
}

func CheckRegexp(field string, value string, reg *regexp.Regexp) error {
	if !reg.MatchString(value) {
		return NewErrInvalidEntity(field, "value does not match required pattern")
	}

	return nil
}

func CheckRestricted[T comparable](field string, source T, allowed ...T) error {
	for _, a := range allowed {
		if source == a {
			return nil
		}
	}

	return NewErrNotAllowed(field, allowed...)
}
