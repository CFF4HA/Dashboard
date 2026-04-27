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
var (
	concurrencyLimiter = make(chan struct{}, 5) // limit to 5 concurrent requests to PubChem
)

// this will create an ingredient given a name, and will contact the PubChem API
// directly. Do not assume that this is controlled via rate-limiting, etc. Call this
// only when you directly want to create an ingredient, and not as part of a batch process.
func RetrieveIngredient(name string) (*types.Ingredient, error) {
	// we need to do rate limiting so that a 503 doesn't tank our product creation
	// form process...
	for {
		select {
		case concurrencyLimiter <- struct{}{}:
			defer func() { <-concurrencyLimiter }()
			is_compound := true
			name = strings.ToLower(strings.TrimSpace(name))
			pubchem_id, err := pubchem.GetCompoundId(name)
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					pubchem_id, err = pubchem.GetSubstanceId(name)
					if err != nil && strings.Contains(err.Error(), "404") {
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
					is_compound = false
				} else {
					return nil, err
				}
			}

			var data []byte
			if is_compound {
				data, err = pubchem.GetCompoundAsJSON(pubchem_id)
			} else {
				data, err = pubchem.GetSubstanceAsJSON(pubchem_id)
				core.Logger.Debug("retrieved substance data from PubChem", "name", name, "data", string(data))
			}
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

			core.Logger.Info("adverse effects data", "data", _adverse_effects.String())
			core.Logger.Info("hazards data", "data", _ghs_hazards.String())
			core.Logger.Info("symptoms data", "data", _signs_symptoms.String())
			core.Logger.Info("california cosmetics data", "data", _california_cosmetics.String())

			var names []string
			name_data := "[" + strings.Trim(_names.String(), "[]") + "]"
			if err := json.Unmarshal([]byte(name_data), &names); err != nil {
				core.Logger.Error("failed to unmarshal names for an ingredient", "name", name, "err", err)
			}

			var hazards []string
			hazard_data := "[" + strings.Trim(_ghs_hazards.String(), "[]") + "]"
			if _ghs_hazards.Exists() {
				if err := json.Unmarshal([]byte(hazard_data), &hazards); err != nil {
					core.Logger.Error("failed to unmarshal hazards for an ingredient", "name", name, "err", err)
				}
			}

			var symptoms []string
			symptom_data := "[" + strings.Trim(_signs_symptoms.String(), "[]") + "]"
			if _signs_symptoms.Exists() {
				if strings.Count(symptom_data, "[") > 1 {
					// there's more than one array, so we need to split and then iterate and join
					arrays := strings.Split(symptom_data, "],[")
					for _, array := range arrays {
						elements := strings.Split(strings.Trim(array, "[]"), ",")
						symptoms = append(symptoms, elements...)
					}
				} else if err := json.Unmarshal([]byte(symptom_data), &symptoms); err != nil {
					core.Logger.Error("failed to unmarshal symptoms for an ingredient", "name", name, "err", err)
				}
			}

			var effects []string
			effect_data := "[" + strings.Trim(_adverse_effects.String(), "[]") + "]"
			if _adverse_effects.Exists() {
				if err := json.Unmarshal([]byte(effect_data), &effects); err != nil {
					core.Logger.Error("failed to unmarshal effects for an ingredient", "name", name, "err", err)
				}
			}

			var regulations []string
			regulation_data := "[" + strings.Trim(_california_cosmetics.String(), "[]") + "]"
			if _california_cosmetics.Exists() {
				if err := json.Unmarshal([]byte(regulation_data), &regulations); err != nil {
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

		default:
			time.Sleep(200 * time.Millisecond) // wait for 200ms before trying again
		}
	}

}
