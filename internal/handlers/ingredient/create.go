package ingredient

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/CFF4HA/Dashboard/pkg/pubchem"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

// this file conatins the functions required to create a new ingredient in the database.

// this will create an ingredient given a name, and will contact the PubChem API
// directly. Do not assume that this is controlled via rate-limiting, etc. Call this
// only when you directly want to create an ingredient, and not as part of a batch process.
func RetrieveIngredient(name string) (*types.Ingredient, error) {
	pubchem_id, err := pubchem.GetCompoundId(name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			ing := &types.Ingredient{
				Model: types.Model{
					Id:      uuid.New(),
					Created: time.Now(),
					Updated: time.Now(),
				},
				PrimaryName: name,
				Failed:      true,
			}

			return ing, core.DB.Create(ing).Error
		}
		return nil, err
	}

	data, err := pubchem.GetCompoundAsJSON(pubchem_id)
	if err != nil {
		return nil, err
	}

	// the following section extracts from the data the relevant information to create
	// an ingredient (i.e., the hazards, effects, regulations, etc.)
	_signs_symptoms := gjson.GetBytes(data, `Record.Section.#(TOCHeading=="Toxicity").Section.#(TOCHeading=="Toxicological Information").Section.#(TOCHeading=="Signs and Symptoms").Information.#.Value.StringWithMarkup.#.String`)
	_adverse_effects := gjson.GetBytes(data, `Record.Section.#(TOCHeading=="Toxicity").Section.#(TOCHeading=="Toxicological Information").Section.#(TOCHeading=="Adverse Effects").Information.#.Value.StringWithMarkup.#.String`)
	_ghs_hazards := gjson.GetBytes(data, `Record.Section.#(TOCHeading=="Safety and Hazards").Section.#(TOCHeading=="Hazards Identification").Section.#(TOCHeading=="GHS Classification").Information.#(Name=="GHS Hazard Statements").Value.StringWithMarkup.#.String`)
	_california_cosmetics := gjson.GetBytes(data, `Record.Section.#(TOCHeading=="Safety and Hazards").Section.#(TOCHeading=="Regulatory Information").Information.#(Name=="California Safe Cosmetics Program (CSCP) Reportable Ingredient").Value.StringWithMarkup.#.String`)
	_names := gjson.GetBytes(data, `Record.Section.#(TOCHeading=="Names and Identifiers").Section.#(TOCHeading=="Synonyms").Section.#(TOCHeading=="Depositor-Supplied Synonyms").Information.#.Value.StringWithMarkup.#.String`)

	var names []string
	if err := json.Unmarshal([]byte(_names.Array()[0].String()), &names); err != nil {
		core.Logger.Error("failed to unmarshal names for an ingredient", "name", name, "err", err)
	}

	var hazards []string
	if _ghs_hazards.Exists() {
		if err := json.Unmarshal([]byte(_ghs_hazards.Array()[0].String()), &hazards); err != nil {
			core.Logger.Error("failed to unmarshal hazards for an ingredient", "name", name, "err", err)
		}
	}

	var symptoms []string
	if _signs_symptoms.Exists() {
		if err := json.Unmarshal([]byte(_signs_symptoms.Array()[0].String()), &symptoms); err != nil {
			core.Logger.Error("failed to unmarshal symptoms for an ingredient", "name", name, "err", err)
		}
	}

	var effects []string
	if _adverse_effects.Exists() {
		if err := json.Unmarshal([]byte(_adverse_effects.Array()[0].String()), &effects); err != nil {
			core.Logger.Error("failed to unmarshal effects for an ingredient", "name", name, "err", err)
		}
	}

	var regulations []string
	if _california_cosmetics.Exists() {
		if err := json.Unmarshal([]byte(_california_cosmetics.Array()[0].String()), &regulations); err != nil {
			core.Logger.Error("failed to unmarshal regulations for an ingredient", "name", name, "err", err)
		}
	}

	ingredient := &types.Ingredient{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},

		PrimaryName: name,
		Failed:      false,
		Metadata: types.IngredientMetadata{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
		},
	}
	ingredient.Metadata.IngredientId = ingredient.Id
	ingredient.Metadata.NumHazards = len(hazards)
	ingredient.Metadata.NumEffects = len(effects)
	ingredient.Metadata.NumSymptoms = len(symptoms)
	ingredient.Metadata.NumRegulations = len(regulations)

	for _, n := range names {
		var existing types.Name
		tx := core.DB.Where("text = ?", n).First(&existing)
		if tx.Error == nil {
			ingredient.Names = append(ingredient.Names, existing)
		} else {
			ingredient.Names = append(ingredient.Names, types.Name{
				Model: types.Model{
					Id:      uuid.New(),
					Created: time.Now(),
					Updated: time.Now(),
				},
				Text: n,
			})
		}
	}

	for _, hazard := range hazards {
		label := types.Label{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
			Type:    "hazard",
			Payload: hazard,
			Origin:  new(string),
		}
		*label.Origin = fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/compound/%d", pubchem_id)
		ingredient.Labels = append(ingredient.Labels, label)
	}

	for _, effect := range effects {
		label := types.Label{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
			Type:    "effect",
			Payload: effect,
			Origin:  new(string),
		}
		*label.Origin = fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/compound/%d", pubchem_id)
		ingredient.Labels = append(ingredient.Labels, label)
	}

	for _, symptom := range symptoms {
		label := types.Label{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
			Type:    "symptom",
			Payload: symptom,
			Origin:  new(string),
		}
		*label.Origin = fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/compound/%d", pubchem_id)
		ingredient.Labels = append(ingredient.Labels, label)
	}

	for _, regulation := range regulations {
		label := types.Label{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
			Type:    "regulation",
			Payload: regulation,
			Origin:  new(string),
		}

		*label.Origin = fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/compound/%d", pubchem_id)
		ingredient.Labels = append(ingredient.Labels, label)
	}

	if tx := core.DB.Create(ingredient); tx.Error != nil {
		core.Logger.Error("failed to create ingredient in database", "name", name, "err", tx.Error)
		return ingredient, tx.Error
	}

	return ingredient, nil
}
