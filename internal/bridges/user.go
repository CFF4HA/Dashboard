package bridges

import (
	"errors"
	"net/http"
	"slices"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

type UserSessionRequired struct {
	Excludes []string
}

func (u UserSessionRequired) Data(w http.ResponseWriter, r *http.Request, model map[string]any) (any, error) {
	if slices.Contains(u.Excludes, r.URL.Path) {
		return nil, nil
	}

	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}

	// if the value exists but is empty, then the user needs to
	// return to the login page to get a new session cookie.
	if c.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, errors.New("session cookie is empty")
	}

	provided_session, err := uuid.Parse(c.Value)
	if err != nil {
		return nil, errors.New("invalid session cookie")
	}

	user_id, ok := user.GetSession(provided_session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, errors.New("invalid session")
	}

	user := &types.User{}
	if err := core.DB.First(user, "id = ?", user_id).Error; err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}

	return user, nil
}

func (u UserSessionRequired) Name() string {
	return "User"
}
