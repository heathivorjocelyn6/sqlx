package sqlx

import (
	"errors"
	"reflect"
	"strings"
)

// In expands a query with ? placeholders to use the correct number of placeholders
// for slice arguments.
func In(query string, args ...interface{}) (string, []interface{}, error) {
	var newArgs []interface{}
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			if v.Len() == 0 {
				return "", nil, errors.New("empty slice passed to 'in' query")
			}
		}
		newArgs = append(newArgs, arg)
	}
	// ... existing logic for query expansion ...
	return query, newArgs, nil
}