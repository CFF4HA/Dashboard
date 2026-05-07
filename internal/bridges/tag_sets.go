package bridges

import (
	"github.com/CFF4HA/Dashboard/internal/backend"
	"net/http"
)

var (
	TagSets = DataBridge{Key: "TagSets", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetTaggingSets(r.FormValue("cursor"))
	}}

	TagSetsByname = DataBridge{Key: "TagSetsByName", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetTaggingSetsByName(r.FormValue("name"), r.FormValue("cursor"))
	}}
)
