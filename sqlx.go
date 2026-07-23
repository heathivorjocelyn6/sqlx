package sqlx

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func In(query string, args ...interface{}) (string, []interface{}, error) {
	var newQuery strings.Builder
	var newArgs []interface{}
	argIdx := 0

	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			if argIdx >= len(args) {
				return "", nil, errors.New("number of placeholders doesn't match number of arguments")
			}
			arg := args[argIdx]
			argIdx++

			v := reflect.ValueOf(arg)
			if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
				if v.Len() == 0 {
					newQuery.WriteString("NULL")
				} else {
					for j := 0; j < v.Len(); j++ {
						if j > 0 {
							newQuery.WriteString(",")
						}
						newQuery.WriteByte('?')
						newArgs = append(newArgs, v.Index(j).Interface())
					}
				}
			} else {
				newQuery.WriteByte('?')
				newArgs = append(newArgs, arg)
			}
		} else {
			newQuery.WriteByte(query[i])
		}
	}

	if argIdx != len(args) {
		return "", nil, errors.New("number of placeholders doesn't match number of arguments")
	}

	return newQuery.String(), newArgs, nil
}

func Rebind(query string) string {
	var b strings.Builder
	n := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			b.WriteByte('$')
			b.WriteString(strconv.Itoa(n))
			n++
		} else {
			b.WriteByte(query[i])
		}
	}
	return b.String()
}
