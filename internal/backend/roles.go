package backend

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func GetRoles() ([]types.Role, error) {
	var roles []types.Role

	if err := core.DB.Order("id").Find(&roles).Error; err != nil {
		core.Logger.Error("failed to retrieve roles", "err", err)
		return nil, errors.New("failed to retrieve roles, try again later.")
	}

	return roles, nil
}

func RouteGetRoles(w http.ResponseWriter, r *http.Request) error {
	roles, err := GetRoles()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(roles)
}
