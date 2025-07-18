package routes

import (
	"github.com/Estate-CRM/backend-go/internal/handlers"
	"github.com/go-chi/chi"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()
	authHandler := &handlers.AuthHandler{}
	r.Post("/login", authHandler.Login)
	r.Post("/registerclient", authHandler.RegisterClient)
	r.Post("/registeragent", authHandler.RegisterAgent)
	r.Post("/testdata", authHandler.Testdata)
	/* r.Post("/login", loginHandler)
	r.Post("/register", registerHandler) */

	return r
}
