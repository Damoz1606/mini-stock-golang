package main

import (
	"log"
	"net/http"
	"os"

	helloinfra "github.com/Damoz1606/ministock-backend/internal/modules/query/hello/infrastructure/transport/http"
	"github.com/Damoz1606/ministock-backend/internal/modules/query/hello/usecase/gethello"
	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	usecase := gethello.NewGetHelloHandler()
	handler := helloinfra.NewHandler(usecase)
	router := helloinfra.NewRouter(handler)

	chiRouter := chi.NewRouter()
	router.Register(chiRouter)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: chiRouter,
	}

	log.Printf("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}