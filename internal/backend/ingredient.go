package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/CFF4HA/Dashboard/pkg/pubchem"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// ------------------------------
// Generic Utility Functions, not Route Related
//
// These functions provide the functionality the routes need to perform
// actions.
// ------------------------------

// This function is the function responsible for pulling ingredient information from the data sources.
func pullIngredientByName(name string) (*types.Ingredient, error) {
	// Will always upload an ingredient, marks it as failed if it did not
	// work as expected. This is to prevent repeated attempts to pull the same ingredient.
	var rawData *json.RawMessage

	name = strings.ToLower(strings.TrimSpace(name))
	ing := &types.Ingredient{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		PrimaryName: name,
	}
	core.Logger.Debug("attempting to pull ingredient information from pubchem", "name", name)

	cid, err := pubchem.GetCompoundId(name)
	if err != nil {
		core.Logger.Error("failed to retrieve compound id from pubchem", "name", name, "error", err)

		sid, err := pubchem.GetSubstanceId(name)
		if err != nil {
			ing.Failed = true
			return ing, core.DB.Create(&ing).Error
		}

		substance, err := pubchem.GetSubstanceAsJSON(sid)
		if err != nil {
			core.Logger.Error("failed to retrieve substance information from pubchem", "name", name, "sid", sid, "error", err)
			ing.Failed = true
			return ing, core.DB.Create(&ing).Error
		}

		rawData := &json.RawMessage{}
		*rawData = substance
	}

	if rawData == nil {
		compound, err := pubchem.GetCompoundAsJSON(cid)
		if err != nil {
			core.Logger.Error("failed to retrieve compound information from pubchem", "name", name, "cid", cid, "error", err)

			ing.Failed = true
			tx := core.DB.Create(&ing)
			return ing, tx.Error
		}

		rawData = &json.RawMessage{}
		*rawData = compound
	}

	// TODO: Parse information from the configurable database backed pubchem config
	// set by the system.
	allLabelConfigs := []types.PubChemLabelConfig{}
	if tx := core.DB.Model(&types.PubChemLabelConfig{}).Find(&allLabelConfigs); tx.Error != nil {
		core.Logger.Error("failed to retrieve pubchem label configs from database", "error", tx.Error)

		ing.Failed = true
		tx := core.DB.Create(&ing)
		return ing, tx.Error
	}

	var labels []types.Label
	for _, cfg := range allLabelConfigs {
		// Finish implementation of the gjson parsing
		// and marshalling them into arrays of strings.
		data := gjson.GetBytes(*rawData, cfg.GsonField)
		core.Logger.Debug("retrieved data from pubchem for label config", "name", name, "config", cfg, "data", data.String())

		// For each string in the array of strings, create a new label and match it
		// with the created ingredient.
		var text []string
		if data.Exists() {
			normalized := "[" + strings.Trim(data.String(), "[]") + "]"
			if err := json.Unmarshal([]byte(normalized), &text); err != nil {
				core.Logger.Warn("failed to unmarshal data from pubchem for label config", "name", name, "config", cfg, "data", data.String(), "error", err)
				continue
			}
		}

		for _, content := range text {
			label := &types.Label{
				Model: types.Model{
					Id:      uuid.New(),
					Created: time.Now(),
					Updated: time.Now(),
				},
				IngredientId: ing.Id,
				Type:         cfg.LabelType,
				Payload:      content,
			}
			label.Origin = new(string)
			*label.Origin = fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/compound/%d", cid)

			labels = append(labels, *label)
		}
	}
	ing.Labels = labels

	// TODO: Do the tagging based on the rules in the system.

	currentPubChemLookupsCountLock.Lock()
	currentPubChemLookupsCount++
	currentPubChemLookupsCountLock.Unlock()

	core.Logger.Info("successfully pulled ingredient information from pubchem", "name", name, "cid", cid, "labels", labels)
	tx := core.DB.Create(&ing)
	if tx.Error != nil {
		core.Logger.Error("failed to insert ingredient into database after pulling from pubchem", "name", name, "error", tx.Error)
		return ing, tx.Error
	}

	return ing, nil
}

func InsertIngredientNote(user_id uuid.UUID, id string, note string) error {
	// we first will retrieve the ingredient, then update the note field,
	// then save it back to the database with the note.

	ing, err := GetIngredientById(id)
	if err != nil {
		core.Logger.Error("failed to retrieve ingredient from database for note insertion", "id", id, "error", err)
		return errors.New("failed to retrieve ingredient from database, try again later.")
	}

	n := &types.IngredientNote{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Content:      note,
		UserId:       user_id,
		IngredientId: ing.Model.Id,
	}

	tx := core.DB.Create(n)
	if tx.Error != nil {
		core.Logger.Error("failed to insert ingredient note into database", "id", id, "error", tx.Error)
		return errors.New("failed to insert ingredient note into database, try again later.")
	}

	return nil
}

func RetrieveIngredientByPrimaryName(name string) (*types.Ingredient, error) {
	name = strings.ToLower(strings.TrimSpace(name))

	ing := &types.Ingredient{}
	tx := core.DB.Where("primary_name = ?", name).First(ing)
	if tx.Error != nil {
		// if the error is a record not found error we need to perform a search for
		// the same chemical using the NAMES table, but it must be an
		// exact match.
		if tx.Error == gorm.ErrRecordNotFound {
			return pullIngredientByName(name)
		}

		core.Logger.Error("failed to retrieve ingredient from database", "name", name, "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredient from database")
	}

	return ing, nil
}

func TagIngredient(ingredient_id string, tag_id string) error {
	tx := core.DB.Exec("INSERT INTO ingredient_tags (ingredient_id, tag_id) VALUES (?, ?)", ingredient_id, tag_id)
	if tx.Error != nil {
		core.Logger.Error("failed to tag ingredient", "ingredient_id", ingredient_id, "tag_id", tag_id, "error", tx.Error)
		return errors.New("failed to tag ingredient, try again later.")
	}

	return nil
}

func RemoveTagFromIngredient(ingredient_id string, tag_id string) error {
	tx := core.DB.Exec("DELETE FROM ingredient_tags WHERE ingredient_id = ? AND tag_id = ?", ingredient_id, tag_id)
	if tx.Error != nil {
		core.Logger.Error("failed to remove tag from ingredient", "ingredient_id", ingredient_id, "tag_id", tag_id, "error", tx.Error)
		return errors.New("failed to remove tag from ingredient, try again later.")
	}

	return nil
}

func DeleteIngredientById(id string) error {
	tx := core.DB.Exec("DELETE FROM ingredient_tags WHERE ingredient_id = ?", id)
	tx = core.DB.Exec("DELETE FROM ingredient_notes WHERE ingredient_id = ?", id)
	tx = core.DB.Exec("DELETE FROM product_ingredients WHERE ingredient_id = ?", id)

	tx = core.DB.Delete(&types.Ingredient{}, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to delete ingredient from database", "id", id, "error", tx.Error)
		return errors.New("failed to delete ingredient from database, try again later.")
	}

	return nil
}

func GetIngredients(cursor string) ([]types.Ingredient, error) {
	var ingredients []types.Ingredient

	tx := core.DB.Scopes(WithCursor(cursor), WithLimit(20), WithOrder("id"), WithPreload("Labels", "Tags")).Find(&ingredients)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredients from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredients from database, try again later.")
	}

	return ingredients, nil
}

func GetIngredientsByPrimaryName(name string, cursor string) ([]types.Ingredient, error) {
	var ingredients []types.Ingredient

	tx := core.DB.Scopes(WithCursor(cursor), WithLimit(20), WithSearch("primary_name", name), WithOrder("id"),
		WithPreload("Labels", "Tags", "Notes"),
	).
		Find(&ingredients)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredients from database", "name", name, "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredients from database, try again later.")
	}

	return ingredients, nil
}

func GetIngredientById(id string) (*types.Ingredient, error) {
	var ingredient *types.Ingredient

	tx := core.DB.Scopes(WithPreload("Labels", "Tags")).First(&ingredient, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredient from database", "id", id, "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredient from database, try again later.")
	}

	return ingredient, nil
}

func GetIngredientNotes(ingredient_id string, u uuid.UUID, cursor string) ([]types.IngredientNote, error) {
	var notes []types.IngredientNote

	tx := core.DB.Scopes(WithCursor(cursor), WithLimit(20), WithOrder("id")).Where("user_id = ?", u).Where("ingredient_id = ?", ingredient_id).Find(&notes)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredient notes from database", "ingredient_id", ingredient_id, "user_id", u, "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredient notes from database, try again later.")
	}

	return notes, nil
}

func DeleteIngredientNote(note_id string, u uuid.UUID) error {
	tx := core.DB.Exec("DELETE FROM ingredient_notes WHERE id = ? AND user_id = ?", note_id, u)
	tx = core.DB.Delete(&types.IngredientNote{}, "id = ?", note_id, "user_id = ?", u)
	if tx.Error != nil {
		core.Logger.Error("failed to delete ingredient note from database", "note_id", note_id, "user_id", u, "error", tx.Error)
		return errors.New("failed to delete ingredient note from database, try again later.")
	}

	return nil
}

// ------------------------------
// Route Related Functions
//
// These routes provide actual route functionality.
// ------------------------------

func RouteRetrieveIngredientByPrimaryName(w http.ResponseWriter, r *http.Request) error {
	ingredient, err := RetrieveIngredientByPrimaryName(r.FormValue("name"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ingredient)
}

func RouteTagIngredient(w http.ResponseWriter, r *http.Request) error {
	return TagIngredient(r.FormValue("ingredient_id"), r.FormValue("tag_id"))
}

func RouteRemoveTagFromIngredient(w http.ResponseWriter, r *http.Request) error {
	return RemoveTagFromIngredient(r.FormValue("ingredient_id"), r.FormValue("tag_id"))
}

func RouteRemoveIngredientById(w http.ResponseWriter, r *http.Request) error {
	return DeleteIngredientById(r.FormValue("id"))
}

func RouteGetIngredients(w http.ResponseWriter, r *http.Request) error {
	ingredients, err := GetIngredients(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ingredients)
}

func RouteGetIngredientsByPrimaryName(w http.ResponseWriter, r *http.Request) error {
	ingredients, err := GetIngredientsByPrimaryName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ingredients)
}

func RouteGetIngredientById(w http.ResponseWriter, r *http.Request) error {
	ing, err := GetIngredientById(r.FormValue("id"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ing)
}

func RouteIngredientAddNote(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return errors.New("a valid user is required to perform this operation, please try again after signing in")
	}

	err = InsertIngredientNote(u.Id, r.FormValue("id"), r.FormValue("note"))
	if err != nil {
		return err
	}

	return nil
}

func RouteDeleteIngredientNote(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return err
	}

	return DeleteIngredientNote(r.FormValue("id"), u.Id)
}

func RouteGetIngredientNotes(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequest(w, r)
	if err != nil || u == nil {
		return errors.New("a valid user is required to perform this operation, please try again after signing in")
	}

	notes, err := GetIngredientNotes(r.FormValue("id"), u.Id, r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(notes)
}
