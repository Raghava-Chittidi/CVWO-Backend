package router

import (
	"github.com/CVWO-Backend/internal/middlewares"
	"github.com/CVWO-Backend/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CORS)
	r.Group(routes.UnrestrictedRoutes())
	r.Group(func(s chi.Router) {
		s.Use(middlewares.AuthoriseUser)
		s.Group(routes.RestrictedRoutes())
	})
}
