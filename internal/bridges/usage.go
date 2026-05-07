package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

var (
	LatestUsageMetrics = DataBridge{Key: "UsageMetric", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetLatestUsageMetric()
	}}
)
