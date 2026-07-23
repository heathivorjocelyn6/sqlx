package sqlx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrEmptySlice is returned when an empty slice is passed to In for expansion.
var ErrEmptySlice = errors.New("empty slice passed to 'in' query")

// In expands slice arguments in a query containing IN (?) placeholders.
//
// For each ? in the query, if the corresponding arg is a slice, the ? is
// replaced with the appropriate number of placeholders (e.g. ?,?,? for a
// 3-element slice). Non-slice args are left as-is.
//
// If any slice argument has length 0, In returns ErrEmptySlice to prevent
// generating invalid SQL like "IN ()".
func In(query string, args ...interface{}) (string, []interface{}, error) {
	// First pass: validate args and collect expanded args
	expandedArgs := make([]interface{}, 0, len(args))

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			if v.Len() == 0 {
				return "", nil, ErrEmptySlice
			}
			// Expand slice into individual args
			for i := 0; i < v.Len(); i++ {
				expandedArgs = append(expandedArgs, v.Index(i).Interface())
			}
		} else {
			expandedArgs = append(expandedArgs, arg)
		}
	}

	// Second pass: replace ? placeholders with expanded placeholders
	var result strings.Builder
	argIndex := 0
	i := 0
	for i < len(query) {
		if query[i] == '?' {
			if argIndex >= len(args) {
				return "", nil, fmt.Errorf("more placeholders than arguments")
			}
			v := reflect.ValueOf(args[argIndex])
			if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
				if v.Len() == 0 {
					return "", nil, ErrEmptySlice
				}
				// Write expanded placeholders (no extra parens — the query
				// already has IN (?) and we replace the ? with ?,?,?)
				for j := 0; j < v.Len(); j++ {
					if j > 0 {
						result.WriteString(",")
					}
					result.WriteString("?")
				}
			} else {
				result.WriteString("?")
			}
			argIndex++
			i++
		} else {
			result.WriteByte(query[i])
			i++
		}
	}

	return result.String(), expandedArgs, nil
}
