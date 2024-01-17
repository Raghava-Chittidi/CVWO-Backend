package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type auth struct {
	Issuer string
	Audience string
	Secret string
	CookiePath string
	CookieName string
	CookieDomain string
	TokenExpiry time.Duration
	RefreshExpiry time.Duration
}

type AuthenticatedUser struct {
	ID int `json:"id"`
	Username string `json:"username"`
}

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var Auth auth

func GenerateAuth() {
	Auth = auth{
		Issuer: "example.com",
		Audience: "example.com",
		Secret: "keyboardsecret",
		CookiePath: "/",
		CookieName: "session-cookie",
		CookieDomain: "",
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
	}
}

func (j *auth) GenerateTokens(user *AuthenticatedUser) (Tokens, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["sub"] = fmt.Sprint(user.ID)
	claims["iss"] = j.Issuer
	claims["aud"] = j.Audience
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	signedAccessToken, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return Tokens{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return Tokens{}, err
	}

	tokens := Tokens {
		AccessToken: signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokens, nil
}

func (j *auth) GenerateRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: refreshToken,
		Domain: j.CookieDomain,
		MaxAge: int(j.RefreshExpiry.Seconds()),
		Expires: time.Now().Add(j.RefreshExpiry),
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		HttpOnly: true,
	}
}

func (j *auth) DeleteRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: "",
		Domain: j.CookieDomain,
		MaxAge: -1,
		Expires: time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		HttpOnly: true,
	}
}

func (j *auth) VerifyToken(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errors.New("No authorization header")
	}

	arr := strings.Split(authHeader, " ")
	if len(arr) != 2 || arr[0] != "Bearer" {
		return "", nil, errors.New("Invalid authorization header")
	}

	token := arr[1]
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC);
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("Invalid issuer")
	}

	return token, claims, nil
}