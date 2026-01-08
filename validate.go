package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Validate takes a struct and validates values of its fields based on their tags.
func Validate(obj interface{}, options *ValidateOptions) (bool, map[string]int, error) {
	if options == nil {
		options = &ValidateOptions{} // use defaults
	}

	tagName := "validation"
	if options.TagName != "" {
		tagName = options.TagName
	}

	val := reflect.ValueOf(obj)

	// If we received a pointer to a pointer, keep dereferencing until we hit a non‑pointer.
	// This also normalizes the case where the caller passed a pointer to a struct.
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return false, nil, errors.New("nil pointer supplied")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false, nil, fmt.Errorf("obj is not a struct or pointer to struct, got %s", val.Kind())
	}

	typ := val.Type()

	numField := typ.NumField()
	fieldViolations := make(map[string]int, numField)
	overallOk := true

	for i := 0; i < numField; i++ {
		structField := typ.Field(i)
		fieldValue := val.Field(i)

		// unexported fields are not validated
		if !structField.IsExported() {
			continue
		}

		// check if only specified field should be checked
		if len(options.RestrictFields) > 0 && !options.RestrictFields[structField.Name] {
			continue
		}

		// check if the field exists in options.OverwriteValues
		if options.OverwriteValues != nil {
			if overwriteVal, exists := options.OverwriteValues[structField.Name]; exists {
				fieldValue = reflect.ValueOf(overwriteVal)
				// check if fieldValue can be assigned to structField
				if !fieldValue.Type().AssignableTo(structField.Type) {
					overallOk = false
					fieldViolations[structField.Name] = FailType
					continue
				}
			}
		}

		ok, viol := ValidateField(structField, fieldValue, tagName)
		if !ok {
			overallOk = false
			fieldViolations[structField.Name] = viol
		}
	}

	return overallOk, fieldViolations, nil
}

// ValidateField takes a reflected struct field, its value and a tagname and validates the values against the requirements in the tag.
func ValidateField(structField reflect.StructField, fieldValue reflect.Value, tagName string) (bool, int) {
	violations := 0

	tag := structField.Tag.Get(tagName)

	isNilPointer, kind, concrete := dereferenceKind(fieldValue)

	tag = strings.TrimSpace(tag)
	if tag == "-" {
		return true, 0
	}

	if tag != "" {
		tokens := strings.Fields(tag)

		// check req first
		for _, token := range tokens {
			name, _ := parseRule(token)
			if name == "req" && isNilPointer {
				violations += FailReq
				break
			}
		}

		// make further checks only if not nil
		if !isNilPointer {
			for _, token := range tokens {
				name, arg := parseRule(token)
				switch name {
				case "email":
					if kind != reflect.String {
						// Wrong kind → treat as a violation (helps catch config errors)
						violations += FailType
						break
					}
					if !emailRegexp.MatchString(concrete.String()) {
						violations += FailEmail
					}

				case "lenmin", "lenmax":
					if kind != reflect.String {
						violations += FailType
						break
					}
					minMax, _ := strconv.Atoi(arg) // arg is guaranteed to be numeric by convention
					strLen := len(concrete.String())
					if name == "lenmin" && strLen < minMax {
						violations += FailLenMin
					}
					if name == "lenmax" && strLen > minMax {
						violations += FailLenMax
					}
				case "valmin", "valmax":
					// First, make sure we are dealing with a numeric kind.
					switch kind {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						limit, _ := strconv.ParseInt(arg, 10, 64)
						val := concrete.Int()
						if name == "valmin" && val < limit {
							violations += FailValMin
						}
						if name == "valmax" && val > limit {
							violations += FailValMax
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
						limit, _ := strconv.ParseUint(arg, 10, 64)
						val := concrete.Uint()
						if name == "valmin" && val < limit {
							violations += FailValMin
						}
						if name == "valmax" && val > limit {
							violations += FailValMax
						}
					case reflect.Float32, reflect.Float64:
						limit, _ := strconv.ParseFloat(arg, 64)
						val := concrete.Float()
						if name == "valmin" && val < limit {
							violations += FailValMin
						}
						if name == "valmax" && val > limit {
							violations += FailValMax
						}
					default:
						// Not a numeric type → record a mismatch.
						violations = +FailType
					}
				default:
				}
			}
		}
	}

	// regexp checked only if the field value is not a nil pointer
	if !isNilPointer {
		regexpTagName := tagName + "_regexp"
		if pattern := structField.Tag.Get(regexpTagName); pattern != "" {
			if kind == reflect.String {
				compiled, err := getCompiledRegexp(pattern)
				if err != nil || !compiled.MatchString(concrete.String()) {
					// configuration error is considered a failed regexp
					violations += FailRegExp
				}
			}
		}
	}

	ok := violations == 0
	return ok, violations
}
