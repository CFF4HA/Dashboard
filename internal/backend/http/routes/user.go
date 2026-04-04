package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

var (
	Sessions = make(map[uuid.UUID]uuid.UUID)
)

func HandleSessionPOST(w http.ResponseWriter, r *http.Request) error {
	var session types.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		core.Logger.Error("tried to decode request body for session validation, but failed", "error", err)
		return err
	}

	_, ok := Sessions[session.Session]
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func HandleUserPOST(w http.ResponseWriter, r *http.Request) error {
	db := database.Database()

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		core.Logger.Error("tried to decode request body for user creation, but failed", "error", err)
		return err
	}

	var user types.User
	if tx := db.Model(&types.User{}).Where("username = ?", request.Username).First(&user); tx.Error != nil {
		core.Logger.Error("tried to find user in database, but failed", "error", tx.Error)
		return tx.Error
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		core.Logger.Info("password for user login is wrong", "username", request.Username)
		return errors.New("wrong password")
	}

	session := uuid.New()
	Sessions[session] = user.Uid

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})

	return nil
}

func HandleUserPUT(w http.ResponseWriter, r *http.Request) error {
	db := database.Database()

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		core.Logger.Error("tried to decode request body for user creation, but failed", "error", err)
		return err
	}

	hash, err := hashPassword(request.Password)
	if err != nil {
		core.Logger.Error("tried to hash password for user creation, but failed", "error", err)
		return err
	}

	user := &types.User{
		Uid:          uuid.New(),
		Username:     request.Username,
		PasswordHash: hash,
	}
	if tx := db.Model(&types.User{}).Create(user); tx.Error != nil {
		core.Logger.Error("tried to create user in database, but failed", "error", tx.Error)
		return err
	}

	session := uuid.New()
	Sessions[session] = user.Uid

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})

	return nil
}
