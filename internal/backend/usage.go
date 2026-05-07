package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

var (
	currentUsageMetric *types.UsageMetric

	currentAutomatedScrapesCount int
	automatedScrapeCountLock     sync.Mutex

	currentPubChemLookupsCount     int
	currentPubChemLookupsCountLock sync.Mutex
)

func GetProductsCount() int64 {
	var count int64
	tx := core.DB.Model(&types.Product{}).Count(&count)
	if tx.Error != nil {
		return 0
	}

	return count
}

func GetIngredientsCount() int64 {
	var count int64
	tx := core.DB.Model(&types.Ingredient{}).Count(&count)
	if tx.Error != nil {
		return 0
	}

	return count
}

func GetLabelsCount() int64 {
	var count int64
	tx := core.DB.Model(&types.PubChemLabelConfig{}).Count(&count)
	if tx.Error != nil {
		return 0
	}

	return count
}

func GetTagsCount() int64 {
	var count int64
	tx := core.DB.Model(&types.Tag{}).Count(&count)
	if tx.Error != nil {
		return 0
	}

	return count
}

func GetTaggingRuleCount() int64 {
	var count int64
	tx := core.DB.Model(&types.TaggingRule{}).Count(&count)
	if tx.Error != nil {
		return 0
	}

	return count
}

func InsertNewUsageMetric() (*types.UsageMetric, error) {
	if currentUsageMetric == nil {
		currentUsageMetric = &types.UsageMetric{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
		}
	}

	currentUsageMetric.ProductsCount = GetProductsCount()
	currentUsageMetric.IngredientsCount = GetIngredientsCount()
	currentUsageMetric.LabelsCount = GetLabelsCount()
	currentUsageMetric.TagsCount = GetTagsCount()
	currentUsageMetric.TaggingRuleCount = GetTaggingRuleCount()
	currentUsageMetric.AutomatedScrapesCount = currentAutomatedScrapesCount
	currentUsageMetric.PubChemLookupsCount = currentPubChemLookupsCount
	currentUsageMetric.SessionsCount = len(user.SessionMap)

	tx := core.DB.Create(currentUsageMetric)
	if tx.Error != nil {
		core.Logger.Error("failed to upload usage metrics", "err", tx.Error)
		return currentUsageMetric, errors.New("failed to insert usage metrics to the database, please try again later")
	}

	return currentUsageMetric, nil
}

func UsageMetricDaemon() {
	ticker := time.NewTicker(10 * time.Minute)
	InsertNewUsageMetric()

	for range ticker.C {
		InsertNewUsageMetric()
	}
}

func GetLatestUsageMetric() (*types.UsageMetric, error) {
	metric := &types.UsageMetric{}

	tx := core.DB.Order("created desc").First(metric)
	if tx.Error != nil {
		return nil, errors.New("failed to get latest usage metric from the database, please try again later")
	}

	return metric, nil
}

// -------------------------------

func RouteGetLatestUsageMetric(w http.ResponseWriter, r *http.Request) error {
	m, err := GetLatestUsageMetric()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(m)
}
