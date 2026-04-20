package user

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	SessionMap = make(map[uuid.UUID]uuid.UUID)
	sessionMu  sync.RWMutex
)

// getSession returns the user ID for a session token, and whether it exists.
func GetSession(session uuid.UUID) (uuid.UUID, bool) {
	sessionMu.RLock()
	defer sessionMu.RUnlock()
	userID, ok := SessionMap[session]
	return userID, ok
}

// setSession maps a session token to a user ID.
func SetSession(session, userID uuid.UUID) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	SessionMap[session] = userID
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// This function will attempt to create a new user. It requires
// a username and password in the request body, sent using form
// encoding.
//
// The password will be hashed and stored in the database.
func HandleUserPUT(w http.ResponseWriter, r *http.Request) error {
	username := strings.TrimSpace(r.FormValue("username"))
	password_input := (r.FormValue("password"))

	if username == "" || password_input == "" {
		return errors.New("username and password are required")
	}

	hash, err := hashPassword(password_input)
	if err != nil {
		return err
	}

	var user types.User
	user.Model.Id = uuid.New()
	user.Model.Created = time.Now()
	user.Model.Updated = time.Now()
	user.Username = username
	user.PasswordHash = hash

	session := uuid.New()
	SetSession(session, user.Model.Id)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	return core.DB.Create(&user).Error
}

// HandleUserPOST handles login. It reads username and password from the form,
// verifies the password against the stored hash, and sets a session cookie.
// An optional `to` query parameter is respected for post-login redirection.
func HandleUserPOST(w http.ResponseWriter, r *http.Request) error {
	username := strings.TrimSpace(r.FormValue("username"))
	password_input := r.FormValue("password")

	if username == "" || password_input == "" {
		return errors.New("username and password are required")
	}

	var u types.User
	if err := core.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password_input)); err != nil {
		return errors.New("invalid username or password")
	}

	session := uuid.New()
	SetSession(session, u.Model.Id)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	to := r.URL.Query().Get("to")
	if to == "" || !strings.HasPrefix(to, "/") {
		to = "/"
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}
