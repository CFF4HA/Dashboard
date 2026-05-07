package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

func InsertNotification(content string, color string, enable bool) (*types.Notification, error) {
	notification := &types.Notification{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Content:         content,
		BackgroundColor: color,
		Enabled:         enable,
	}

	tx := core.DB.Create(notification)
	if tx.Error != nil {
		core.Logger.Error("failed to insert notification into database", "error", tx.Error)
		return nil, errors.New("failed to insert notification into database, try again later")
	}

	core.Logger.Debug("successfully inserted notification into database", "id", notification.Id)
	return notification, nil
}

func DeleteNotificationById(id string) error {
	tx := core.DB.Delete(&types.Notification{}, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to delete notification from database", "error", tx.Error, "id", id)
		return errors.New("failed to delete notification from database, try again later")
	}

	return nil
}

func GetNotificationsByEnabled(enabled bool, cursor string) ([]types.Notification, error) {
	var notifications []types.Notification

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithLimit(20)).Where("enabled = ?", enabled).Find(&notifications)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve notifications from database", "error", tx.Error, "enabled", enabled)
		return nil, errors.New("failed to retrieve notifications from database, try again later")
	}

	return notifications, nil
}

func GetNotifications(cursor string) ([]types.Notification, error) {
	var notifications []types.Notification

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithLimit(20)).Find(&notifications)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve notifications from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve notifications from database, try again later")
	}

	return notifications, nil
}

// -------------------------------

func RouteInsertNotification(w http.ResponseWriter, r *http.Request) error {
	core.Logger.Debug("received request to insert notification", "content", r.FormValue("content"), "color", r.FormValue("color"), "enabled", r.FormValue("enabled"))

	notification, err := InsertNotification(r.FormValue("content"), r.FormValue("color"), r.FormValue("enabled") == "true")
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(notification)
}

func RouteDeleteNotification(w http.ResponseWriter, r *http.Request) error {
	return DeleteNotificationById(r.FormValue("id"))
}

func RouteGetNotifications(w http.ResponseWriter, r *http.Request) error {
	notifications, err := GetNotifications(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(notifications)
}

func RouteGetNotificationsByEnabled(w http.ResponseWriter, r *http.Request) error {
	notifications, err := GetNotificationsByEnabled(r.FormValue("enabled") == "true", r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(notifications)
}
