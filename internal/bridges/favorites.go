package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/google/uuid"
)

var UserFavoriteProductsBridge = verb.Map("FavoriteProducts", func(r *http.Request, m map[string]any) (any, error) {
	u, ok := m["User"].(*types.User)
	if !ok || u == nil {
		return map[uuid.UUID]types.Product{}, nil
	}

	var favs []types.Product
	if err := core.DB.Model(u).Association("Products").Find(&favs); err != nil {
		return map[uuid.UUID]types.Product{}, nil
	}

	if len(favs) > 0 {
		var ids []uuid.UUID
		for _, p := range favs {
			ids = append(ids, p.Model.Id)
		}
		core.DB.Preload("Ingredients").Preload("Metadata").Where("id IN ?", ids).Find(&favs)
	}

	result := make(map[uuid.UUID]types.Product, len(favs))
	for _, p := range favs {
		result[p.Model.Id] = p
	}
	return result, nil
})

var UserFavoriteIngredientsBridge = verb.Map("FavoriteIngredients", func(r *http.Request, m map[string]any) (any, error) {
	u, ok := m["User"].(*types.User)
	if !ok || u == nil {
		return map[uuid.UUID]types.Ingredient{}, nil
	}

	var favs []types.Ingredient
	if err := core.DB.Model(u).Association("Ingredients").Find(&favs); err != nil {
		return map[uuid.UUID]types.Ingredient{}, nil
	}

	if len(favs) > 0 {
		var ids []uuid.UUID
		for _, p := range favs {
			ids = append(ids, p.Model.Id)
		}
		core.DB.Preload("Metadata").Preload("Names").Preload("Labels").Where("id IN ?", ids).Find(&favs)
	}

	result := make(map[uuid.UUID]types.Ingredient, len(favs))
	for _, i := range favs {
		result[i.Model.Id] = i
	}
	return result, nil
})
