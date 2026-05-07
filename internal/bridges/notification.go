package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

var (
	Notifications = DataBridge{Key: "Notifications", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetNotifications(r.FormValue("cursor"))
	}}

	NotificationsEnabled = DataBridge{Key: "NotificationsEnabled", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetNotificationsByEnabled(r.FormValue("enabled") == "true", r.FormValue("cursor"))
	}}
)
