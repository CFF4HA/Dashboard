package core

import (
	"net/http"
	"strings"
)

func RequestValuesAsSlice(r *http.Request, key string) []string {
	values, ok := r.Form[key]
	if !ok {
		return []string{}
	}

	elements := []string{}
	for _, v := range values {
		for _, elem := range strings.Split(v, ",") {
			value := strings.TrimSpace(elem)
			if len(value) > 0 {
				elements = append(elements, elem)
			}
		}
	}

	return elements
}
