package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/Estate-CRM/backend-go/internal/utils"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) *utils.AppError

// ErrorHandler is the middleware that catches errors and returns HTTP response
func ErrorHandler(fn AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
		}
	}
}
