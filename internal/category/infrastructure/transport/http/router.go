package http

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", handler.Create)
	r.Get("/", handler.ReadAll)
	r.Get("/{categoryId}", handler.ReadOne)
	r.Put("/{categoryId}", handler.Update)
	r.Delete("/{categoryId}", handler.Delete)

	return r
}
