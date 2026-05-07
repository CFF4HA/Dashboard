package bridges

import (
	"net/http"
)

type DataBridge struct {
	Key  string
	Func func(w http.ResponseWriter, r *http.Request) (any, error)
}

func (b DataBridge) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	return b.Func(w, r)
}

func (b DataBridge) Name() string {
	return b.Key
}

// ------------------
