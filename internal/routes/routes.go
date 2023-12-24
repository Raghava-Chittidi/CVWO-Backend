package routes

import (
	handlers "github.com/CVWO-Backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/create/category", handlers.CreateCategory)
		r.Get("/categories", handlers.GetCategories)
		r.Get("/threads", handlers.GetThreads)
		r.Get("/threads/{id}", handlers.GetThread)
		r.Get("/refresh", handlers.RefreshToken)
		r.Post("/signup", handlers.SignUp)
		r.Post("/login", handlers.Authenticate)
		r.Post("/logout", handlers.Logout)
	}
}

func RestrictedRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/create/thread", handlers.CreateThread)
		r.Post("/create/comment", handlers.CreateComment)
		r.Patch("/edit/comment/{id}", handlers.EditComment)
		r.Delete("/delete/comment/{id}", handlers.DeleteComment)
	}
}


