package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/CVWO-Backend/internal/auth"
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"github.com/CVWO-Backend/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type authInfo struct {
	AccessToken string
	RefreshToken string
	Email string
	Username string
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == auth.Auth.CookieName {

			claims := &auth.Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(auth.Auth.Secret), nil
			})
			if err != nil {
				util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
			}

			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				util.ErrorJSON(w, errors.New("Unknown user!"), http.StatusUnauthorized)
			}

			user, err := GetUserById(userId)
			if err != nil {
				util.ErrorJSON(w, errors.New("Unknown user!"), http.StatusUnauthorized)
			}

			u := auth.JwtUser{
				ID: int(user.ID),
				Username: user.Username,
			}

			tokenPair, err := auth.Auth.GenerateTokenPair(&u)
			if err != nil {
				util.ErrorJSON(w, errors.New("Error generating tokens!"), http.StatusUnauthorized)
			}

			http.SetCookie(w, auth.Auth.GenerateRefreshCookie(tokenPair.RefreshToken))
			authInfo := authInfo{AccessToken: tokenPair.Token, RefreshToken: tokenPair.RefreshToken, Email: user.Email, Username: user.Username}
			util.WriteJSON(w, authInfo, http.StatusOK)
		}
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	// Read JSON payload
	var reqPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &reqPayload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	// Validate user against database
	user, err := GetUserByEmail(reqPayload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("Invalid Credentials!"))
		return
	}

	// Check password
	valid, err := user.VerifyPassword(reqPayload.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("Invalid Credentials!"))
		return
	}

	// Create a jwtUser
	jwtUser := auth.JwtUser {
		ID: int(user.ID),
		Username: user.Username,
	}

	// Generate tokens
	tokens, err := auth.Auth.GenerateTokenPair(&jwtUser)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}
	authInfo := authInfo{AccessToken: tokens.Token, RefreshToken: tokens.RefreshToken, Email: user.Email, Username: user.Username}
	refreshCookie := auth.Auth.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, authInfo, http.StatusAccepted)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string
		Email string
		Password string
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Username == "" {
		util.ErrorJSON(w, errors.New("Invalid username!"))
		return
	}

	_, err = GetUserByUsername(payload.Username)
	if err != gorm.ErrRecordNotFound {
		util.ErrorJSON(w, errors.New("Username is in use!"))
		return
	}

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, errors.New("Invalid email!"))
		return
	} 

	_, err = GetUserByEmail(payload.Email)
	if err != gorm.ErrRecordNotFound {
		util.ErrorJSON(w, errors.New("Email is in use!"))
		return
	}
	
	if payload.Password == "" {
		util.ErrorJSON(w, errors.New("Invalid password!"))
		return
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	user := models.User{Username: payload.Username, Email: payload.Email, Password: string(hashedPw)}

	result := database.DB.Create(&user)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	// Create a jwtUser
	jwtUser := auth.JwtUser {
		ID: int(user.ID),
		Username: user.Username,
	}

	// Generate tokens
	tokens, err := auth.Auth.GenerateTokenPair(&jwtUser)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}
	authInfo := authInfo{AccessToken: tokens.Token, RefreshToken: tokens.RefreshToken, Email: user.Email, Username: user.Username}
	refreshCookie := auth.Auth.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, authInfo, http.StatusAccepted)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, auth.Auth.DeleteRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUserById(id int) (*models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}