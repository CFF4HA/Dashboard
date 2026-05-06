package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

func InsertPubChemLabelConfig(gjson string, label_type string) (*types.PubChemLabelConfig, error) {
	cfg := &types.PubChemLabelConfig{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		GsonField: gjson,
		LabelType: label_type,
	}

	tx := core.DB.Create(cfg)
	if tx.Error != nil {
		core.Logger.Error("failed to insert pubchem label config", "error", tx.Error)
		return nil, errors.New("failed to insert pubchem label config, please try again later.")
	}

	return cfg, nil
}

func DeletePubChemLabelConfig(id uuid.UUID) error {
	tx := core.DB.Delete(&types.PubChemLabelConfig{}, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to delete pubchem label config", "error", tx.Error)
		return errors.New("failed to delete pubchem label config, please try again later.")
	}

	return nil
}

func GetPubChemLabelConfigs(cursor string) ([]types.PubChemLabelConfig, error) {
	var cfgs []types.PubChemLabelConfig
	tx := core.DB.Scopes(
		WithCursor(cursor),
		WithLimit(20),
		WithOrder("id"),
	).Find(&cfgs)
	if tx.Error != nil {
		core.Logger.Error("failed to get pubchem label configs", "error", tx.Error)
		return nil, errors.New("failed to get pubchem label configs, please try again later.")
	}

	return cfgs, nil
}

// --------------------------

func RouteInsertPubChemLabelConfig(w http.ResponseWriter, r *http.Request) error {
	cfg, err := InsertPubChemLabelConfig(r.FormValue("gjson"), r.FormValue("label_type"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(cfg)
}

func RouteDeletePubChemLabelConfig(w http.ResponseWriter, r *http.Request) error {
	return DeletePubChemLabelConfig(uuid.MustParse(r.FormValue("id")))
}

func RouteGetPubChemLabelConfigs(w http.ResponseWriter, r *http.Request) error {
	cfgs, err := GetPubChemLabelConfigs(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(cfgs)
}
