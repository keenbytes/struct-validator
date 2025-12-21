package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// dereferenceKind walks through pointer indirections until it reaches a non‑pointer type.
// It returns the final reflect.Kind and the reflect.Value that points to that concrete value.
// If the original value is a nil pointer, the returned value will be the zero Value of the
// element type (so Kind() will be the element’s kind, but IsValid() will be false).
func dereferenceKind(v reflect.Value) (bool, reflect.Kind, reflect.Value) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			// Nil pointer – we cannot Elem() safely, so break and return the zero value.
			// The caller can decide whether a nil pointer is acceptable.
			return true, v.Type().Elem().Kind(), reflect.Zero(v.Type().Elem())
		}
		v = v.Elem()
	}
	return false, v.Kind(), v
}

func parseRule(tok string) (name, arg string) {
	parts := strings.SplitN(tok, ":", 2)
	name = parts[0]
	if len(parts) == 2 {
		arg = parts[1]
	}
	return
}

var (
	// pattern → compiled *regexp.Regexp*
	regexCache sync.Map
)

// getCompiledRegexp returns a compiled *regexp.Regexp for the given pattern.
// It caches the result so subsequent calls are O(1).
func getCompiledRegexp(pattern string) (*regexp.Regexp, error) {
	if v, ok := regexCache.Load(pattern); ok {
		return v.(*regexp.Regexp), nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexCache.Store(pattern, compiled)
	return compiled, nil
}
