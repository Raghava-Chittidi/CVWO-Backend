package routes

import (
	handlers "github.com/CVWO-Backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/createCategory", handlers.CreateCategory)
		r.Get("/categories", handlers.GetCategories)
		r.Get("/threads", handlers.GetThreads)
		r.Get("/refresh", handlers.RefreshToken)
		r.Post("/signup", handlers.SignUp)
		r.Post("/login", handlers.Authenticate)
		r.Post("/createThread", handlers.CreateThread)
		r.Get("/logout", handlers.Logout)

		r.Get("/", handlers.Home)
	}
}


