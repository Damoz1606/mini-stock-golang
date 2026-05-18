package helloinfra

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	handler http.Handler
}

func NewRouter(handler http.Handler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Register(router chi.Router) {
	router.Get("/v1/hello", r.handler.ServeHTTP)
}