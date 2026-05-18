package helloinfra

import (
	"encoding/json"
	"net/http"

	"github.com/Damoz1606/ministock-backend/internal/modules/query/hello/usecase/gethello"
)

type Handler struct {
	usecase *gethello.GetHelloHandler
}

func NewHandler(usecase *gethello.GetHelloHandler) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	query := gethello.GetHelloQuery{}
	response, err := h.usecase.Handle(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}