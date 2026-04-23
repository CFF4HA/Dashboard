package bridges

import (
	"net/http"
	"slices"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

type InvestigationBridge struct {
}

func (i InvestigationBridge) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	// we need to get the list of products (comma separated) from the
	// key called "compare_product_ids"
	product_ids_str := strings.Split(r.Form.Get("compare_product_ids"), ",")

	// we now need to convert these to actual product types from the database,
	// and they need to be fully populated.
	var products []types.Product
	for _, product_id := range product_ids_str {
		id, err := uuid.Parse(product_id)
		if err != nil {
			core.Logger.Warn("invalid product id was received, skipping", "id", product_id)
			continue
		}

		var product types.Product
		if err := core.DB.Preload("Ingredients").Preload("Metadata").First(&product, "id = ?", id).Error; err != nil {
			core.Logger.Warn("error retrieving product from database, skipping", "id", product_id, "error", err)
			continue
		}

		for i := range product.Ingredients {
			core.DB.Preload("Labels").Preload("Names").Preload("Metadata").First(&product.Ingredients[i])
		}

		products = append(products, product)
	}

	inv := types.Investigation{
		ProductList: products,
	}
	inv.Shared.Ingredients = getSharedIngredients(products)
	inv.Shared.Hazards = getSharedLabels(products, "hazard")
	inv.Shared.Effects = getSharedLabels(products, "effect")
	inv.Shared.Symptoms = getSharedLabels(products, "symptom")
	inv.Shared.Regulations = getSharedLabels(products, "regulation")
	inv.Unique = make(map[uuid.UUID]struct {
		Ingredients []types.Ingredient
		Hazards     []string
		Effects     []string
		Symptoms    []string
		Regulations []string
	})

	// find the unique ingredients, hazards, effects, symptoms, and regulations for each product
	for _, product := range products {
		var unique_ingredients []types.Ingredient
		var unique_hazards []string
		var unique_effects []string
		var unique_symptoms []string
		var unique_regulations []string

		for _, ingredient := range product.Ingredients {
			found := false
			for _, sharedIngredient := range inv.Shared.Ingredients {
				if ingredient.Id == sharedIngredient.Id {
					found = true
					break
				}
			}

			for _, label := range ingredient.Labels {
				switch label.Type {
				case "hazard":
					if !slices.Contains(inv.Shared.Hazards, label.Payload) {
						unique_hazards = append(unique_hazards, label.Payload)
					}
				case "effect":
					if !slices.Contains(inv.Shared.Effects, label.Payload) {
						unique_effects = append(unique_effects, label.Payload)
					}
				case "symptom":
					if !slices.Contains(inv.Shared.Symptoms, label.Payload) {
						unique_symptoms = append(unique_symptoms, label.Payload)
					}
				case "regulation":
					if !slices.Contains(inv.Shared.Regulations, label.Payload) {
						unique_regulations = append(unique_regulations, label.Payload)
					}
				default:
					continue
				}
			}

			if found {
				continue
			}

			unique_ingredients = append(unique_ingredients, ingredient)
		}

		inv.Unique[product.Id] = struct {
			Ingredients []types.Ingredient
			Hazards     []string
			Effects     []string
			Symptoms    []string
			Regulations []string
		}{
			Ingredients: unique_ingredients,
			Hazards:     unique_hazards,
			Effects:     unique_effects,
			Symptoms:    unique_symptoms,
			Regulations: unique_regulations,
		}
	}

	return inv, nil
}

func (i InvestigationBridge) Name() string {
	return "Investigation"
}

// -- Helper Functions --
func getSharedIngredients(products []types.Product) []types.Ingredient {
	ingredients := []types.Ingredient{}
	existence_map := make(map[uuid.UUID]int)
	ingredient_map := make(map[uuid.UUID]types.Ingredient)

	for _, product := range products {
		for _, ingredient := range product.Ingredients {
			_, ok := existence_map[ingredient.Id]
			if !ok {
				existence_map[ingredient.Id] = 1
				ingredient_map[ingredient.Id] = ingredient
			} else {
				existence_map[ingredient.Id] += 1
			}
		}
	}

	for k, v := range existence_map {
		if v > 1 {
			ingredients = append(ingredients, ingredient_map[k])
		}
	}

	return ingredients
}

func getSharedLabels(products []types.Product, labelType string) []string {
	labels := []string{}
	existence_map := make(map[string]int)
	label_map := make(map[string]string)

	for _, product := range products {
		for _, ingredient := range product.Ingredients {
			for _, label := range ingredient.Labels {
				if label.Type == labelType {
					if _, ok := existence_map[label.Payload]; !ok {
						existence_map[label.Payload] = 1
						label_map[label.Payload] = label.Payload
					} else {
						existence_map[label.Payload] += 1
					}
				}
			}
		}
	}

	for k, v := range existence_map {
		if v > 1 {
			labels = append(labels, label_map[k])
		}
	}

	return labels
}
