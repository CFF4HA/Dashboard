package database

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

type ingredientNameRequest struct {
	Name        string `json:"name"`
	CacheBypass bool   `json:"cache_bypass"`
}

var (
	ingredientUpdateChannel                              = make(chan ingredientNameRequest, 20)
	ingredientCache         map[string]*types.Ingredient = make(map[string]*types.Ingredient)
	ingredientCacheLock                                  = &sync.RWMutex{}
)

func SyncDatabaseWithIngredient(name string, cacheBypass bool) {
	// this will do a request to the python backend that should trigger a
	// scrape of PubChem
	ingredientCacheLock.RLock()
	v, ok := ingredientCache[name]
	if ok && !cacheBypass {
		core.Logger.Info("ingredient cache hit", "name", name)
		ingredientCacheLock.RUnlock()
		return
	}
	ingredientCacheLock.RUnlock()

	ingredientCacheLock.Lock()
	defer ingredientCacheLock.Unlock()

	v, ok = ingredientCache[name]
	if ok && v != nil && !cacheBypass {
		core.Logger.Info("ingredient cache hit", "name", name)
		return
	}

	core.Logger.Info("syncing ingredient with database", "name", name, "cache_bypass", cacheBypass)
	request, err := http.NewRequest("GET", core.BackendAddress+"/ingredient?name="+url.QueryEscape(name), nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	ingredientCache[name] = nil
}

func IngredientUpdateMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case request := <-ingredientUpdateChannel:
			SyncDatabaseWithIngredient(request.Name, request.CacheBypass)
		}
	}
}

func PollForIngredient(name string, timeout_s int, query_interval_s int, cache_bypass bool) (*types.Ingredient, error) {
	var ingredient *types.Ingredient
	n := new(types.Name)

	found := false
	timeout := false

	// setup the infrastructure to query for the ingredient.
	ticker := time.NewTicker(time.Duration(query_interval_s) * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout_s)*time.Second)
	defer cancel()

	if cache_bypass {
		SyncDatabaseWithIngredient(name, cache_bypass)
	}

	// Ping the database, or timeout.
	for !found && !timeout {
		select {
		case <-ticker.C:
			ingredientCacheLock.RLock()

			v, ok := ingredientCache[name]
			if ok && v != nil && !cache_bypass {
				core.Logger.Info("ingredient cache hit", "name", name)
				ingredientCacheLock.RUnlock()
				return v, nil
			}
			ingredientCacheLock.RUnlock()

			if tx := db.Model(&types.Name{}).Where("text ~* ?", name).Preload("Ingredient").First(&n); tx.Error != nil {
				SyncDatabaseWithIngredient(name, cache_bypass)
				continue
			}

			ingredient = &n.Ingredient
			db.Preload("Labels").Preload("Names").Find(ingredient)
			found = true

			core.Logger.Info("ingredient found in database", "name", name)
		case <-ctx.Done():
			core.Logger.Debug("ingredient search timed out", "name", name)
			timeout = true
		}
	}

	if !found {
		return nil, errors.New("ingredient not found")
	}

	return ingredient, nil
}
