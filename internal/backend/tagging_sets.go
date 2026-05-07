package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// -------------------------------
// Generic Utility Functions, not Route Related
// --------------------------------

func InsertTaggingSet(name string, description string) (*types.TaggingSet, error) {
	set := &types.TaggingSet{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Name:        strings.ToLower(strings.TrimSpace(name)),
		Description: description,
	}

	tx := core.DB.Create(set)
	if tx.Error != nil {
		core.Logger.Error("failed to insert tagging set into database", "error", tx.Error)
		return nil, errors.New("failed to insert tagging set into database, try again later")
	}

	core.Logger.Debug("successfully inserted tagging set into database", "tagging_set_id", set.Id)
	return set, nil
}

func DeleteTaggingSet(setId string) error {
	tx := core.DB.Delete(&types.TaggingSet{}, "id = ?", setId)
	if tx.Error != nil {
		core.Logger.Error("failed to delete tagging set from database", "tagging_set_id", setId, "error", tx.Error)
		return errors.New("failed to delete tagging set from database, try again later")
	}

	core.Logger.Debug("successfully deleted tagging set from database", "tagging_set_id", setId)
	return nil
}

func GetTaggingSets(cursor string) ([]types.TaggingSet, error) {
	var sets []types.TaggingSet

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithPreload("Rules", "Users")).Find(&sets)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging sets from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging sets from database, try again later")
	}

	for i := range sets {
		for j := range sets[i].Rules {
			core.DB.Model(&sets[i].Rules[j]).Association("Tag").Find(&sets[i].Rules[j].Tag)
		}
	}

	core.Logger.Debug("successfully retrieved tagging sets from database", "count", len(sets))
	return sets, nil
}

func GetTaggingSetsByName(name string, cursor string) ([]types.TaggingSet, error) {
	var sets []types.TaggingSet

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithSearch("name", name), WithPreload("Rules", "Users")).Find(&sets)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging sets from database", "name", name, "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging sets from database, try again later")
	}

	return sets, nil
}

func GetTaggingSetById(setId string) (*types.TaggingSet, error) {
	var set *types.TaggingSet

	if tx := core.DB.Scopes(WithPreload("Rules", "Users")).First(&set, "id = ?", setId); tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging set from database", "tagging_set_id", setId, "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging set from database, try again later")
	}

	return set, nil
}

// --------------------------------
// Route Handlers
// --------------------------------

func RouteInsertTaggingSet(w http.ResponseWriter, r *http.Request) error {
	set, err := InsertTaggingSet(r.FormValue("name"), r.FormValue("description"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(set)
}

func RouteDeleteTaggingSet(w http.ResponseWriter, r *http.Request) error {
	return DeleteTaggingSet(r.FormValue("id"))
}

func RouteGetTaggingSets(w http.ResponseWriter, r *http.Request) error {
	sets, err := GetTaggingSets(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(sets)
}

func RouteGetTaggingSetById(w http.ResponseWriter, r *http.Request) error {
	set, err := GetTaggingSetById(r.FormValue("id"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(set)
}

func RouteGetTaggingSetByName(w http.ResponseWriter, r *http.Request) error {
	sets, err := GetTaggingSetsByName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(sets)
}
