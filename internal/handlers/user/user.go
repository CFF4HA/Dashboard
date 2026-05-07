package user

import (
	"errors"
	"fmt"
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

// This is a convenience function that will attempt to retrieve the user associated with
// the request. It throws an error on many potential failure points.
func GetUserFromRequestNoRedirect(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	// request must have session cookie
	c, err := r.Cookie("session")
	if err != nil {
		core.Logger.Warn("request missing session cookie")
		return nil, err
	}

	// session cookie must be a valid uuid
	session_id, err := uuid.Parse(c.Value)
	if err != nil {
		core.Logger.Warn("request session cookie was not valid uuid")
		return nil, err
	}

	// uuid must be a valid session in the session map
	user_uid, ok := SessionMap[session_id]
	if !ok {
		core.Logger.Warn("user session cookie was not in session map")
		return nil, err
	}

	// the returned user uid must exist in the database
	user := &types.User{}
	user.Model.Id = user_uid
	if core.DB.First(user).Error != nil {
		core.Logger.Warn("user session cookie had invalid user id not in the database")
		return nil, err
	}

	return user, nil
}

// This is a convenience function that will attempt to retrieve the user associated with
// the request. It throws an error on many potential failure points.
func GetUserFromRequest(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	// request must have session cookie
	c, err := r.Cookie("session")
	if err != nil {
		core.Logger.Warn("request missing session cookie")
		return nil, err
	}

	// session cookie must be a valid uuid
	session_id, err := uuid.Parse(c.Value)
	if err != nil {
		core.Logger.Warn("request session cookie was not valid uuid")
		return nil, err
	}

	// uuid must be a valid session in the session map
	user_uid, ok := SessionMap[session_id]
	if !ok {
		core.Logger.Warn("user session cookie was not in session map")
		return nil, errors.New("user session not in session map")
	}

	// the returned user uid must exist in the database
	user := &types.User{}
	user.Model.Id = user_uid
	if core.DB.First(user).Error != nil {
		core.Logger.Warn("user session cookie had invalid user id not in the database")
		return nil, err
	}

	return user, nil
}

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

	w.Header().Set("Hx-Redirect", "/")
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

func sessionUserID(r *http.Request) (uuid.UUID, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return uuid.Nil, errors.New("not authenticated")
	}
	sessionID, err := uuid.Parse(c.Value)
	if err != nil {
		return uuid.Nil, errors.New("invalid session cookie")
	}
	userID, ok := GetSession(sessionID)
	if !ok {
		return uuid.Nil, errors.New("invalid session")
	}
	return userID, nil
}

func HandleAddUserProduct(w http.ResponseWriter, r *http.Request) error {
	userID, err := sessionUserID(r)
	if err != nil {
		return err
	}

	productID, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		return errors.New("invalid product id")
	}

	u := types.User{}
	u.Model.Id = userID

	var count int64
	core.DB.Table("user_products").
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count)

	if count > 0 {
		var product types.Product
		product.Model.Id = productID
		core.DB.Model(&u).Association("Products").Delete(&product)
		return writeStarButton(w, "/user/product", productID.String(), false)
	}

	var product types.Product
	if err := core.DB.First(&product, "id = ?", productID).Error; err != nil {
		return err
	}
	core.DB.Model(&u).Association("Products").Append(&product)
	return writeStarButton(w, "/user/product", productID.String(), true)
}

func HandleAddUserIngredient(w http.ResponseWriter, r *http.Request) error {
	userID, err := sessionUserID(r)
	if err != nil {
		return err
	}

	ingredientID, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		return errors.New("invalid ingredient id")
	}

	u := types.User{}
	u.Model.Id = userID

	var count int64
	core.DB.Table("user_ingredients").
		Where("user_id = ? AND ingredient_id = ?", userID, ingredientID).
		Count(&count)

	if count > 0 {
		var ingredient types.Ingredient
		ingredient.Model.Id = ingredientID
		core.DB.Model(&u).Association("Ingredients").Delete(&ingredient)
		return writeStarButton(w, "/user/ingredient", ingredientID.String(), false)
	}

	var ingredient types.Ingredient
	if err := core.DB.First(&ingredient, "id = ?", ingredientID).Error; err != nil {
		return err
	}
	core.DB.Model(&u).Association("Ingredients").Append(&ingredient)
	return writeStarButton(w, "/user/ingredient", ingredientID.String(), true)
}

func writeStarButton(w http.ResponseWriter, path, id string, starred bool) error {
	w.Header().Set("Content-Type", "text/html")
	if starred {
		fmt.Fprintf(w, `<button class="btn-star starred" hx-post="%s" hx-vals='{"id": "%s"}' hx-swap="outerHTML" title="Starred"><i class="bi bi-star-fill"></i></button>`, path, id)
	} else {
		fmt.Fprintf(w, `<button class="btn-star" hx-post="%s" hx-vals='{"id": "%s"}' hx-swap="outerHTML" title="Save"><i class="bi bi-star"></i></button>`, path, id)
	}
	return nil
}
