package routes

import (
	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", AuthRoutes())
		r.Mount("/contact",ContactRoutes())
		r.Mount("/property",PropertyRoutes())
	})

	return r
}
