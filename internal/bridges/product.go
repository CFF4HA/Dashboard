package bridges

import (
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/google/uuid"
)

type ProductsByName struct {
}

func (p ProductsByName) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	u, _ := user.GetUserFromRequestNoRedirect(w, r)
	var user_id *uuid.UUID
	if u != nil {
		user_id = &u.Model.Id
	}

	products, err := backend.GetProductsByName(r.FormValue("name"), r.FormValue("cursor"), user_id)
	if err != nil {
		return nil, errors.New("failed to get products by name: " + err.Error())
	}

	core.Logger.Debug("searched for products by name", "name", r.FormValue("name"), "cursor", r.FormValue("cursor"), "user_id", user_id, "results_count", len(products))
	return products, nil
}

func (p ProductsByName) Name() string {
	return "Products"
}
