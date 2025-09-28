package database

import (
	"fmt"
	"strings"
	"time"
)

func DebugQueryWithArgs(queryName string, queryString string, args []interface{}) string {
	query := queryString
	for _, arg := range args {
		val := ""
		switch v := arg.(type) {
		case string:
			val = v
		case []byte:
			val = string(v)
		case time.Time:
			val = v.Format(time.RFC3339)
		default:
			val = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", val, 1)
	}
	return fmt.Sprintf("\n%s\n%s\n", queryName, query)
}

func ComposedQuery(queryString string, args []interface{}) string {
	query := queryString
	for _, arg := range args {
		val := ""
		switch v := arg.(type) {
		case string:
			val = v
		case []byte:
			val = string(v)
		case time.Time:
			val = v.Format(time.RFC3339)
		default:
			val = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", val, 1)
	}
	return query
}
