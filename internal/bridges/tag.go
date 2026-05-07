package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

var (
	Tags = DataBridge{Key: "Tags", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetTags(r.FormValue("cursor"))
	}}

	TagsByName = DataBridge{Key: "Tags", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetTagsByName(r.FormValue("name"), r.FormValue("cursor"))
	}}
)
