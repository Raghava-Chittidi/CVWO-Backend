package middlewares

import (
	"log"
	"net/http"
	"os"

	"github.com/CVWO-Backend/internal/auth"
	"github.com/joho/godotenv"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		godotenv.Load("../.env")
		clientUrl := os.Getenv("CLIENT_URL")
		w.Header().Set("Access-Control-Allow-Origin", clientUrl)
		w.Header().Set("Access-Control-Allow-Credentials", "true")


		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func AuthoriseUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := auth.Auth.VerifyToken(w, r)
		log.Println(err)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}