package routes

import (
	"github.com/Estate-CRM/backend-go/internal/handlers"
	"github.com/go-chi/chi"
)

func ContactRoutes() chi.Router {
	r :=chi.NewRouter()
	r.Post("/create", handlers.CreateContact)
	r.Delete("/delete", handlers.DeleteContact)
	r.Get("/getAll", handlers.GetContacts)

	return r
}