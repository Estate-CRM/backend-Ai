package routes

import (
	"github.com/Estate-CRM/backend-go/internal/handlers"
	"github.com/go-chi/chi"
)

func MatchRoutes() chi.Router {
	r := chi.NewRouter()

	matchHandler := &handlers.MatchHandler{}
	r.Post("/generateContract", matchHandler.HandleGenerateContract)
	r.Post("/createMatch", handlers.CreateMatch)

	return r
}
