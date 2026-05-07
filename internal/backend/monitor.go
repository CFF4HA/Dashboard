package backend

import (
	"errors"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

func InsertMonitor(ingredient_id string, user_id string) (*types.Monitor, error) {
	var ingredient types.Ingredient
	if tx := core.DB.First(&ingredient, "id = ?", ingredient_id); tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredient for monitor insertion", "ingredient_id", ingredient_id, "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredient, try again later.")
	}

	monitor := &types.Monitor{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		IngredientId:          ingredient_id,
		UserId:                user_id,
		LastKnownUpdatedField: ingredient.Updated,
	}

	if tx := core.DB.Create(monitor); tx.Error != nil {
		core.Logger.Error("failed to insert monitor into database", "error", tx.Error)
		return nil, errors.New("failed to insert monitor into database, try again later.")
	}

	core.Logger.Debug("successfully inserted monitor", "ingredient_id", ingredient_id, "user_id", user_id)
	return monitor, nil
}

func CheckIngredientUpdates(user_id string) ([]types.Monitor, error) {
	var monitors []types.Monitor

	tx := core.DB.Table("monitors").
		Select("monitors.*").
		Joins("JOIN ingredients ON CAST(ingredients.id AS varchar) = monitors.ingredient_id").
		Where("monitors.user_id = ?", user_id).
		Where("ingredients.updated > monitors.last_known_updated_field").
		Find(&monitors)

	if tx.Error != nil {
		core.Logger.Error("failed to check monitor updates", "user_id", user_id, "error", tx.Error)
		return nil, errors.New("failed to check monitor updates, try again later.")
	}

	core.Logger.Debug("successfully checked monitor updates", "user_id", user_id, "count", len(monitors))
	return monitors, nil
}

func GetMonitors(user_id string) ([]types.Monitor, error) {
	var monitors []types.Monitor

	tx := core.DB.Where("user_id = ?", user_id).Find(&monitors)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve monitors", "user_id", user_id, "error", tx.Error)
		return nil, errors.New("failed to retrieve monitors, try again later.")
	}

	core.Logger.Debug("successfully retrieved monitors", "user_id", user_id, "count", len(monitors))
	return monitors, nil
}
