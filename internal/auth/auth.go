package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type auth struct {
	Issuer string
	Audience string
	Secret string
	TokenExpiry time.Duration
	RefreshExpiry time.Duration
	CookieDomain string
	CookiePath string
	CookieName string
}

type JwtUser struct {
	ID int `json:"id"`
	Username string `json:"username"`
}

type TokenPair struct {
	Token string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

// type authInfo struct {
// 	domain string
// 	auth auth
// 	JWTSecret string
// 	JWTIssuer string
// 	JWTAudience string
// 	CookieDomain string
// }

var Auth auth

func GenerateAuthInfo() {
	Auth = auth{
		Issuer: "example.com",
		Audience: "example.com",
		Secret: "keyboardsecret",
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath: "/",
		CookieDomain: "",
		CookieName: "session-cookie",
	}
}

func (j *auth) GenerateTokenPair(user *JwtUser) (TokenPair, error) {
	// Create a new access token
	accessToken := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// Set the expiry for access token
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create signed access token
	signedAccessToken, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPair{}, err
	}

	// Create a new refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed access token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPair{}, err
	}

	tokenPair := TokenPair {
		Token: signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPair, nil
}

func (j *auth) GenerateRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: refreshToken,
		Expires: time.Now().Add(j.RefreshExpiry),
		MaxAge: int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteNoneMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}

func (j *auth) DeleteRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: "",
		Expires: time.Unix(0, 0),
		MaxAge: -1,
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}