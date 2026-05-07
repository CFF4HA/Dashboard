package bridges

import (
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

var (
	ProductById = DataBridge{Key: "Product", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
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

	ProductComparison = DataBridge{Key: "Investigation", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		raw := r.FormValue("compare_product_ids")
		if strings.TrimSpace(raw) == "" {
			return &types.ProductComparison{}, nil
		}

		uuids := []uuid.UUID{}
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := uuid.Parse(part)
			if err != nil {
				continue
			}
			uuids = append(uuids, id)
		}

		return backend.CompareProductsByIdList(uuids)
	}}
)
