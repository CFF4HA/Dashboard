package bridges

import (
	"github.com/CFF4HA/Dashboard/internal/backend"
	"net/http"
)

var (
	PubChemLabelConfigs = DataBridge{Key: "PubChemLabelConfigs", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetPubChemLabelConfigs(r.FormValue("cursor"))
	}}
)
