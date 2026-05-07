package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/google/uuid"
)

var (
	ProductById = DataBridge{Key: "ProductById", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		u, _ := user.GetUserFromRequestNoRedirect(w, r)
		var user_id *uuid.UUID
		if u != nil {
			user_id = &u.Model.Id
		}

		return backend.GetProductById(r.FormValue("id"), user_id)
	}}

	ProductsByName = DataBridge{Key: "Products", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		u, _ := user.GetUserFromRequestNoRedirect(w, r)
		var user_id *uuid.UUID
		if u != nil {
			user_id = &u.Model.Id
		}

		return backend.GetProductsByName(r.FormValue("name"), r.FormValue("cursor"), user_id)
	}}
)
