package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

var (
	TaggingRules = DataBridge{Key: "TaggingRules", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetTaggingRules(r.FormValue("cursor"))
	}}
)
