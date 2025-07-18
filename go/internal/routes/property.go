package routes

import (
	"github.com/Estate-CRM/backend-go/internal/handlers"
	"github.com/go-chi/chi"
)

func PropertyRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/create", handlers.CreateProperty)
	r.Delete("/delete", handlers.DeleteProperty)
	r.Get("/getAll", handlers.GetProperties)

	return r
}
