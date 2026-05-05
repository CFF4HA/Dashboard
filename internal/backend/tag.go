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

// ------------------------------
// Generic Utility Functions, not Route Related
//
// The following are utility functions that do not assume the
// method the data was retrieved with.
// ------------------------------

func InsertTag(name string, description string) (*types.Tag, error) {
	tag := &types.Tag{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Name:        strings.ToLower(strings.TrimSpace(name)),
		Description: description,
	}

	// now we will actually insert it into the database
	tx := core.DB.Create(tag)
	if tx.Error != nil {
		core.Logger.Error("failed to insert tag into database", "error", tx.Error)
		return nil, errors.New("failed to insert tag into database, try again later.")
	}

	core.Logger.Debug("successfully inserted tag into database", "name", name, "id", tag.Id)
	return tag, nil
}

func DeleteTagById(id string) error {
	tx := core.DB.Delete(&types.Tag{}, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to delete tag from database", "error", tx.Error, "id", id)
		return errors.New("failed to delete tag from database, try again later.")
	}

	core.Logger.Debug("successfully deleted tag from database", "id", id)
	return nil
}

func GetTags(cursor string) ([]types.Tag, error) {
	var tags []types.Tag

	tx := core.DB.Scopes(WithCursor(cursor), WithLimit(20), WithOrder("id")).Find(&tags)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tags from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve tags from database, try again later.")
	}

	core.Logger.Debug("successfully retrieved tags from database", "count", len(tags), "tags", tags)
	return tags, nil
}

func GetTagsByName(name string, cursor string) ([]types.Tag, error) {
	var tags []types.Tag

	tx := core.DB.Scopes(WithCursor(cursor), WithLimit(20), WithSearch("name", name), WithOrder("id")).Find(&tags)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tags from database", "error", tx.Error, "name", name)
		return nil, errors.New("failed to retrieve tags from database, try again later.")
	}

	core.Logger.Debug("successfully retrieved tags from database", "count", len(tags), "name", name, "tags", tags)
	return tags, nil
}

func GetTagById(id string) (*types.Tag, error) {
	tag := &types.Tag{}

	tx := core.DB.First(tag, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tag from database", "error", tx.Error, "id", id)
		return nil, errors.New("failed to retrieve tag from database, try again later.")
	}

	core.Logger.Debug("successfully retrieved tag from database", "id", id, "tag", tag)
	return tag, nil
}

// ------------------------------
// Routing Related Functions
//
// These are the actual route handlers that the
// router will use.
// ------------------------------
func RouteTagCreate(w http.ResponseWriter, r *http.Request) error {
	tag, err := InsertTag(r.FormValue("name"), r.FormValue("description"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully created tag", "id", tag.Id, "name", tag.Name)
	return nil
}

func RouteTagDelete(w http.ResponseWriter, r *http.Request) error {
	return DeleteTagById(r.FormValue("id"))
}

func RouteTagGet(w http.ResponseWriter, r *http.Request) error {
	tags, err := GetTags(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(tags)
}

func RouteTagGetByName(w http.ResponseWriter, r *http.Request) error {
	tags, err := GetTagsByName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(tags)
}

func RouteTagGetById(w http.ResponseWriter, r *http.Request) error {
	tag, err := GetTagById(r.FormValue("id"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(tag)
}
