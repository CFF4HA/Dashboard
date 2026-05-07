package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func UserLogin(username string, password string) (*uuid.UUID, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	var u types.User
	if err := core.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	session := uuid.New()
	user.SetSession(session, u.Model.Id)
	return &session, nil
}

func InsertUser(username string, password string) (uuid.UUID, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return uuid.Nil, errors.New("username and password are required")
	}

	hash, err := hashPassword(password)
	if err != nil {
		return uuid.Nil, err
	}

	u := types.User{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Username:     username,
		PasswordHash: hash,
	}

	if err := core.DB.Create(&u).Error; err != nil {
		return uuid.Nil, errors.New("failed to create user, try again later.")
	}

	return u.Model.Id, nil
}

func DeleteUser(u uuid.UUID) error {
	if err := core.DB.Delete(&types.User{}, "id = ?", u).Error; err != nil {
		return errors.New("failed to delete user, try again later.")
	}
	return nil
}

func GetUserInformation(u uuid.UUID) (*types.User, error) {
	var user *types.User

	tx := core.DB.Scopes(WithPreload("Roles", "Ingredients", "Products")).First(user, "id = ?", u)
	if tx.Error != nil {
		return nil, errors.New("failed to retrieve user information, try again later.")
	}

	return user, nil
}

func GetUsers(cursor int, count int) ([]uuid.UUID, error) {
	var users []types.User
	if err := core.DB.Offset(cursor).Limit(count).Find(&users).Error; err != nil {
		return nil, errors.New("failed to retrieve users, try again later.")
	}

	ids := make([]uuid.UUID, len(users))
	for i, u := range users {
		ids[i] = u.Model.Id
	}
	return ids, nil
}

func UserRoleAdd(user_id uuid.UUID, role_id string) error {
	var role types.Role
	if err := core.DB.First(&role, "id = ?", role_id).Error; err != nil {
		return errors.New("role not found")
	}

	var u types.User
	if err := core.DB.First(&u, "id = ?", user_id).Error; err != nil {
		return errors.New("user not found")
	}

	if err := core.DB.Model(&role).Association("Users").Append(&u); err != nil {
		return errors.New("failed to add role to user, try again later.")
	}
	return nil
}

func UserRoleRemove(user_id uuid.UUID, role_id string) error {
	var role types.Role
	if err := core.DB.First(&role, "id = ?", role_id).Error; err != nil {
		return errors.New("role not found")
	}

	var u types.User
	if err := core.DB.First(&u, "id = ?", user_id).Error; err != nil {
		return errors.New("user not found")
	}

	if err := core.DB.Model(&role).Association("Users").Delete(&u); err != nil {
		return errors.New("failed to remove role from user, try again later.")
	}
	return nil
}

// ------------------------------
// Routing Related Functions
// ------------------------------

func RouteUserLogin(w http.ResponseWriter, r *http.Request) error {
	session, err := UserLogin(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	return json.NewEncoder(w).Encode(session)
}

func RouteInsertUser(w http.ResponseWriter, r *http.Request) error {
	id, err := InsertUser(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		return err
	}

	session := uuid.New()
	user.SetSession(session, id)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	return json.NewEncoder(w).Encode(id)
}

func RouteDeleteUser(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequestNoRedirect(w, r)
	if err != nil {
		return errors.New("not authenticated")
	}
	return DeleteUser(u.Model.Id)
}

func RouteGetUsers(w http.ResponseWriter, r *http.Request) error {
	cursor, _ := strconv.Atoi(r.FormValue("cursor"))
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil || count <= 0 {
		count = 20
	}

	ids, err := GetUsers(cursor, count)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ids)
}

func RouteUserRoleAdd(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequestNoRedirect(w, r)
	if err != nil {
		return errors.New("not authenticated")
	}
	return UserRoleAdd(u.Model.Id, r.FormValue("role_id"))
}

func RouteUserRoleRemove(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequestNoRedirect(w, r)
	if err != nil {
		return errors.New("not authenticated")
	}
	return UserRoleRemove(u.Model.Id, r.FormValue("role_id"))
}

func RouteGetUserInformation(w http.ResponseWriter, r *http.Request) error {
	u, err := user.GetUserFromRequestNoRedirect(w, r)
	if err != nil || u == nil {
		return errors.New("not authenticated")
	}

	info, err := GetUserInformation(u.Model.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(info)
}
