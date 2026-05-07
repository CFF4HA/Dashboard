package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

var (
	AllUsers = DataBridge{Key: "Users", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetAllUsers()
	}}

	AllRoles = DataBridge{Key: "Roles", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetRoles()
	}}
)
