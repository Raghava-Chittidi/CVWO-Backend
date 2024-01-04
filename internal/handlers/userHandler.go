package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/CVWO-Backend/internal/auth"
	data "github.com/CVWO-Backend/internal/dataaccess"
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"github.com/CVWO-Backend/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type authInfo struct {
	Email string
	Username string
	ID int
	AccessToken string
	RefreshToken string
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
				return
			}

			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				util.ErrorJSON(w, errors.New("Unknown user!"), http.StatusUnauthorized)
				return
			}

			user, err := data.GetUserById(userId)
			if err != nil {
				util.ErrorJSON(w, errors.New("Unknown user!"), http.StatusUnauthorized)
				return
			}

			authenticatedUser := auth.AuthenticatedUser{
				ID: int(user.ID),
				Username: user.Username,
			}

			tokens, err := auth.Auth.GenerateTokens(&authenticatedUser)
			if err != nil {
				util.ErrorJSON(w, errors.New("Error generating tokens!"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, auth.Auth.GenerateRefreshCookie(tokens.RefreshToken))
			authInfo := authInfo{Email: user.Email, Username: user.Username, ID: int(user.ID), 
								 AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}
			util.WriteJSON(w, authInfo, http.StatusOK)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	user, err := data.GetUserByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("Invalid Credentials!"))
		return
	}

	valid, err := user.VerifyPassword(payload.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("Invalid Credentials!"))
		return
	}

	authenticatedUser := auth.AuthenticatedUser {
		ID: int(user.ID),
		Username: user.Username,
	}

	tokens, err := auth.Auth.GenerateTokens(&authenticatedUser)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}
	authInfo := authInfo{Email: user.Email, Username: user.Username, ID: int(user.ID), 
						 AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}
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

	_, err = data.GetUserByUsername(payload.Username)
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

	_, err = data.GetUserByEmail(payload.Email)
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
	result := database.DB.Table("users").Create(&user)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	authenticatedUser := auth.AuthenticatedUser {
		ID: int(user.ID),
		Username: user.Username,
	}

	tokens, err := auth.Auth.GenerateTokens(&authenticatedUser)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}
	authInfo := authInfo{Email: user.Email, Username: user.Username, ID: int(user.ID), 
						 AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}
	refreshCookie := auth.Auth.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, authInfo, http.StatusAccepted)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, auth.Auth.DeleteRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}