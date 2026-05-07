package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
)

var (
	UserInformation = DataBridge{Key: "UserInformation", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		u, err := user.GetUserFromRequestNoRedirect(w, r)
		if err != nil || u == nil {
			return nil, err
		}

		return backend.GetUserInformation(u.Id)
	}}
)
